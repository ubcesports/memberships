package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/utils"
)

var (
	ErrInvalidUserID           = errors.New("invalid user id")
	ErrInvalidTierCode         = errors.New("invalid tier code")
	ErrMembershipNotFound      = errors.New("membership not found")
	ErrMembershipAlreadyExists = errors.New("user already has an active membership")
	ErrTierPriceNotFound       = errors.New("tier price not found")
	ErrStripePriceMissing      = errors.New("stripe price id is not configured")
	ErrStripeNotConfigured     = errors.New("stripe is not configured")
	ErrInvalidStripeSignature  = errors.New("invalid stripe signature")
	ErrInvalidStripeEvent      = errors.New("invalid stripe event")
)

type MembershipService struct {
	membershipRepo *repository.MembershipRepository
	httpClient     *http.Client
}

type stripeCheckoutSession struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type stripeEvent struct {
	Type string `json:"type"`
	Data struct {
		Object stripeCheckoutSessionObject `json:"object"`
	} `json:"data"`
}

type stripeCheckoutSessionObject struct {
	ID            string            `json:"id"`
	PaymentIntent string            `json:"payment_intent"`
	Metadata      map[string]string `json:"metadata"`
}

func NewMembershipService(membershipRepo *repository.MembershipRepository) *MembershipService {
	return &MembershipService{
		membershipRepo: membershipRepo,
		httpClient:     http.DefaultClient,
	}
}

func (s *MembershipService) ListAvailableTiers(ctx context.Context, userID string) (dto.ListMembershipTiersResponse, error) {
	userUUID, err := utils.ParseUUID(userID)
	if err != nil {
		return dto.ListMembershipTiersResponse{}, ErrInvalidUserID
	}

	pricing, err := s.resolvePricingContext(ctx, userUUID)
	if err != nil {
		return dto.ListMembershipTiersResponse{}, err
	}

	rows, err := s.membershipRepo.ListActiveMembershipTiersWithPrices(ctx, db.ListActiveMembershipTiersWithPricesParams{
		Group:         pricing.group,
		StudentStatus: pricing.studentStatus,
	})
	if err != nil {
		return dto.ListMembershipTiersResponse{}, err
	}

	tiers := make([]dto.MembershipTierDTO, 0, len(rows))
	for _, row := range rows {
		tiers = append(tiers, dto.MembershipTierDTO{
			Code:             dto.TierCodeType(row.Code),
			Title:            row.Title,
			Description:      utils.TextPtr(row.Description),
			Price:            utils.NumericString(row.Price),
			Currency:         "CAD",
			Group:            dto.GroupType(row.Group),
			StudentStatus:    dto.StudentStatusType(row.StudentStatus),
			RequiresCheckout: !utils.NumericIsZero(row.Price),
		})
	}

	return dto.ListMembershipTiersResponse{Tiers: tiers}, nil
}

func (s *MembershipService) GetCurrentMembership(ctx context.Context, userID string) (dto.CurrentMembershipResponse, error) {
	userUUID, err := utils.ParseUUID(userID)
	if err != nil {
		return dto.CurrentMembershipResponse{}, ErrInvalidUserID
	}

	row, err := s.membershipRepo.GetCurrentMembershipByUserID(ctx, userUUID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.CurrentMembershipResponse{Membership: nil}, nil
	}
	if err != nil {
		return dto.CurrentMembershipResponse{}, err
	}

	membership := currentMembershipDTO(row)
	return dto.CurrentMembershipResponse{Membership: &membership}, nil
}

func (s *MembershipService) StartCheckout(ctx context.Context, userID string, tierCode dto.TierCodeType) (dto.StartCheckoutResponse, error) {
	userUUID, err := utils.ParseUUID(userID)
	if err != nil {
		return dto.StartCheckoutResponse{}, ErrInvalidUserID
	}
	code, err := parseTierCode(tierCode)
	if err != nil {
		return dto.StartCheckoutResponse{}, err
	}

	if _, err := s.membershipRepo.GetCurrentMembershipByUserID(ctx, userUUID); err == nil {
		return dto.StartCheckoutResponse{}, ErrMembershipAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return dto.StartCheckoutResponse{}, err
	}

	pricing, err := s.resolvePricingContext(ctx, userUUID)
	if err != nil {
		return dto.StartCheckoutResponse{}, err
	}

	tierPrice, err := s.membershipRepo.GetTierPriceByCode(ctx, db.GetTierPriceByCodeParams{
		Code:          code,
		Group:         pricing.group,
		StudentStatus: pricing.studentStatus,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.StartCheckoutResponse{}, ErrTierPriceNotFound
	}
	if err != nil {
		return dto.StartCheckoutResponse{}, err
	}

	now := time.Now()
	startedAt := pgtype.Timestamptz{Time: now, Valid: true}
	expiresAt := pgtype.Timestamptz{Time: expiryForTier(code, now), Valid: true}

	if utils.NumericIsZero(tierPrice.Price) {
		result, err := s.membershipRepo.CreateCompletedPurchase(ctx, repository.CompletedPurchaseParams{
			UserID:                  userUUID,
			TierID:                  tierPrice.TierID,
			GroupAtPurchase:         pricing.group,
			StudentStatusAtPurchase: pricing.studentStatus,
			StartedAt:               startedAt,
			ExpiresAt:               expiresAt,
			StripePaymentIntentID:   pgtype.Text{},
			PriceAmount:             tierPrice.Price,
		})
		if err != nil {
			return dto.StartCheckoutResponse{}, err
		}
		membershipID := utils.UUIDToString(result.Membership.ID)
		return dto.StartCheckoutResponse{
			Status:       "completed",
			MembershipID: &membershipID,
		}, nil
	}

	if !tierPrice.StripePriceID.Valid || tierPrice.StripePriceID.String == "" {
		return dto.StartCheckoutResponse{}, ErrStripePriceMissing
	}

	session, err := s.createStripeCheckoutSession(ctx, stripeCheckoutRequest{
		StripePriceID: tierPrice.StripePriceID.String,
		UserID:        userID,
		TierID:        utils.UUIDToString(tierPrice.TierID),
		TierCode:      string(tierPrice.Code),
		Group:         string(pricing.group),
		StudentStatus: string(pricing.studentStatus),
	})
	if err != nil {
		return dto.StartCheckoutResponse{}, err
	}

	return dto.StartCheckoutResponse{
		Status:                  "checkout_required",
		CheckoutURL:             &session.URL,
		StripeCheckoutSessionID: &session.ID,
	}, nil
}

func (s *MembershipService) HandleStripeWebhook(ctx context.Context, body []byte, signature string) error {
	if err := verifyStripeSignature(body, signature, os.Getenv("STRIPE_WEBHOOK_SECRET")); err != nil {
		return err
	}

	var event stripeEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return ErrInvalidStripeEvent
	}

	if event.Type != "checkout.session.completed" {
		return nil
	}

	session := event.Data.Object
	if session.PaymentIntent == "" {
		return ErrInvalidStripeEvent
	}

	if _, err := s.membershipRepo.GetTransactionByStripePaymentIntentID(ctx, pgtype.Text{String: session.PaymentIntent, Valid: true}); err == nil {
		return nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	userUUID, err := utils.ParseUUID(session.Metadata["user_id"])
	if err != nil {
		return ErrInvalidStripeEvent
	}

	if _, err := s.membershipRepo.GetCurrentMembershipByUserID(ctx, userUUID); err == nil {
		return nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	code, err := parseTierCode(dto.TierCodeType(session.Metadata["tier_code"]))
	if err != nil {
		return ErrInvalidStripeEvent
	}
	group := db.GroupType(session.Metadata["group_at_purchase"])
	studentStatus := db.StudentStatusType(session.Metadata["student_status_at_purchase"])

	tierPrice, err := s.membershipRepo.GetTierPriceByCode(ctx, db.GetTierPriceByCodeParams{
		Code:          code,
		Group:         group,
		StudentStatus: studentStatus,
	})
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = s.membershipRepo.CreateCompletedPurchase(ctx, repository.CompletedPurchaseParams{
		UserID:                  userUUID,
		TierID:                  tierPrice.TierID,
		GroupAtPurchase:         group,
		StudentStatusAtPurchase: studentStatus,
		StartedAt:               pgtype.Timestamptz{Time: now, Valid: true},
		ExpiresAt:               pgtype.Timestamptz{Time: expiryForTier(code, now), Valid: true},
		StripePaymentIntentID:   pgtype.Text{String: session.PaymentIntent, Valid: true},
		PriceAmount:             tierPrice.Price,
	})
	return err
}

type pricingContext struct {
	group         db.GroupType
	studentStatus db.StudentStatusType
}

func (s *MembershipService) resolvePricingContext(ctx context.Context, userID pgtype.UUID) (pricingContext, error) {
	user, err := s.membershipRepo.GetUserForMembershipPricing(ctx, userID)
	if err != nil {
		return pricingContext{}, err
	}

	groups, err := s.membershipRepo.ListUserGroups(ctx, userID)
	if err != nil {
		return pricingContext{}, err
	}

	studentStatus := db.StudentStatusTypeNonStudent
	if user.IsStudent {
		studentStatus = db.StudentStatusTypeStudent
	}

	return pricingContext{
		group:         effectiveGroup(groups),
		studentStatus: studentStatus,
	}, nil
}

func effectiveGroup(groups []db.GroupType) db.GroupType {
	present := make(map[db.GroupType]bool, len(groups))
	for _, group := range groups {
		present[group] = true
	}

	for _, group := range []db.GroupType{
		db.GroupTypeCompetitiveTeam,
		db.GroupTypeBoard,
		db.GroupTypeDirector,
		db.GroupTypeExecutive,
		db.GroupTypeMember,
	} {
		if present[group] {
			return group
		}
	}

	return db.GroupTypeMember
}

func parseTierCode(code dto.TierCodeType) (db.TierCodeType, error) {
	switch db.TierCodeType(code) {
	case db.TierCodeTypeRegular, db.TierCodeTypePremium, db.TierCodeTypeCab, db.TierCodeTypeDay:
		return db.TierCodeType(code), nil
	default:
		return "", ErrInvalidTierCode
	}
}

func expiryForTier(code db.TierCodeType, now time.Time) time.Time {
	location, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		location = time.Local
	}
	localNow := now.In(location)
	if code == db.TierCodeTypeDay {
		return time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), location)
	}

	expiry := time.Date(localNow.Year(), time.April, 30, 23, 59, 59, int(time.Second-time.Nanosecond), location)
	if localNow.After(expiry) {
		expiry = time.Date(localNow.Year()+1, time.April, 30, 23, 59, 59, int(time.Second-time.Nanosecond), location)
	}
	return expiry
}

type stripeCheckoutRequest struct {
	StripePriceID string
	UserID        string
	TierID        string
	TierCode      string
	Group         string
	StudentStatus string
}

func (s *MembershipService) createStripeCheckoutSession(ctx context.Context, arg stripeCheckoutRequest) (stripeCheckoutSession, error) {
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	frontendBaseURL := os.Getenv("FRONTEND_URL")
	if secretKey == "" || frontendBaseURL == "" {
		return stripeCheckoutSession{}, ErrStripeNotConfigured
	}

	form := url.Values{}
	form.Set("mode", "payment")
	form.Set("line_items[0][price]", arg.StripePriceID)
	form.Set("line_items[0][quantity]", "1")
	form.Set("automatic_tax[enabled]", stripeAutomaticTaxEnabled())
	form.Set("success_url", strings.TrimRight(frontendBaseURL, "/")+"/membership/success?session_id={CHECKOUT_SESSION_ID}")
	form.Set("cancel_url", strings.TrimRight(frontendBaseURL, "/")+"/membership")
	form.Set("client_reference_id", arg.UserID)
	form.Set("metadata[user_id]", arg.UserID)
	form.Set("metadata[tier_id]", arg.TierID)
	form.Set("metadata[tier_code]", arg.TierCode)
	form.Set("metadata[group_at_purchase]", arg.Group)
	form.Set("metadata[student_status_at_purchase]", arg.StudentStatus)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.stripe.com/v1/checkout/sessions", strings.NewReader(form.Encode()))
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	req.SetBasicAuth(secretKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return stripeCheckoutSession{}, fmt.Errorf("stripe checkout session failed: status %d: %s", resp.StatusCode, string(data))
	}

	var session stripeCheckoutSession
	if err := json.Unmarshal(data, &session); err != nil {
		return stripeCheckoutSession{}, err
	}
	if session.ID == "" || session.URL == "" {
		return stripeCheckoutSession{}, ErrInvalidStripeEvent
	}
	return session, nil
}

func stripeAutomaticTaxEnabled() string {
	if strings.EqualFold(os.Getenv("STRIPE_AUTOMATIC_TAX_ENABLED"), "true") {
		return "true"
	}
	return "false"
}

func verifyStripeSignature(body []byte, signatureHeader, secret string) error {
	if secret == "" {
		return ErrStripeNotConfigured
	}

	var timestamp, signature string
	for _, part := range strings.Split(signatureHeader, ",") {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) != 2 {
			continue
		}
		switch keyValue[0] {
		case "t":
			timestamp = keyValue[1]
		case "v1":
			signature = keyValue[1]
		}
	}

	if timestamp == "" || signature == "" {
		return ErrInvalidStripeSignature
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return ErrInvalidStripeSignature
	}

	return nil
}

func currentMembershipDTO(row db.GetCurrentMembershipByUserIDRow) dto.CurrentMembershipDTO {
	var transactionStatus *dto.TransactionStatusType
	if row.TransactionStatus.Valid {
		status := dto.TransactionStatusType(row.TransactionStatus.TransactionStatusType)
		transactionStatus = &status
	}

	return dto.CurrentMembershipDTO{
		ID:                      utils.UUIDToString(row.MembershipID),
		TierCode:                dto.TierCodeType(row.Code),
		TierTitle:               row.Title,
		TierDescription:         utils.TextPtr(row.Description),
		GroupAtPurchase:         dto.GroupType(row.GroupAtPurchase),
		StudentStatusAtPurchase: dto.StudentStatusType(row.StudentStatusAtPurchase),
		StartedAt:               row.StartedAt.Time,
		ExpiresAt:               row.ExpiresAt.Time,
		CancelledAt:             utils.TimestamptzPtr(row.CancelledAt),
		Price:                   utils.NumericStringPtr(row.PriceAmount),
		Currency:                "CAD",
		TransactionStatus:       transactionStatus,
		StripePaymentIntentID:   utils.TextPtr(row.StripePaymentIntentID),
	}
}
