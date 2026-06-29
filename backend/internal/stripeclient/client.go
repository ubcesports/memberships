package stripeclient

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/stripe/stripe-go/v86"
)

/*
	Consts and structs
*/

const CheckoutSessionLifetime = time.Hour

type CheckoutSessionRequest struct {
	TransactionID string
	UserID        string
	CustomerEmail string
	PriceID       string
	ProductID     string
	AmountInCents int64 // Amount of money needed to be paid in cents. Needed to determine new price in case user is upgrading their membership
	Currency      string
	IsUpgrade     bool // true if user is upgrading from regular membership -> premium membership
}

type Client struct {
	api        *stripe.Client
	successUrl string
	cancelUrl  string
}

/*
	Public functions
*/

func NewClient() (*Client, error) {
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	successUrl := os.Getenv("STRIPE_CHECKOUT_SUCCESS_URL")
	cancelUrl := os.Getenv("STRIPE_CHECKOUT_CANCEL_URL")

	if secretKey == "" || successUrl == "" || cancelUrl == "" {
		return nil, errors.New("STRIPE_SECRET_KEY, STRIPE_CHECKOUT_SUCCESS_URL, and STRIPE_CHECKOUT_CANCEL_URL are required.")
	}

	return &Client{
		api:        stripe.NewClient(secretKey),
		successUrl: successUrl,
		cancelUrl:  cancelUrl,
	}, nil
}

func (c *Client) GetPrice(ctx context.Context, priceId string) (*stripe.Price, error) {
	return c.api.V1Prices.Retrieve(ctx, priceId, nil)
}

func (c *Client) CreateCheckoutSession(ctx context.Context, request CheckoutSessionRequest) (*stripe.CheckoutSession, error) {
	params := buildCheckoutSessionParams(request, c.successUrl, c.cancelUrl, time.Now())
	params.SetIdempotencyKey(request.TransactionID)
	return c.api.V1CheckoutSessions.Create(ctx, params)
}

func (c *Client) GetCheckoutSession(ctx context.Context, sessionID string) (*stripe.CheckoutSession, error) {
	return c.api.V1CheckoutSessions.Retrieve(ctx, sessionID, nil)
}

func (c *Client) ExpireCheckoutSession(ctx context.Context, sessionID string) (*stripe.CheckoutSession, error) {
	return c.api.V1CheckoutSessions.Expire(ctx, sessionID, nil)
}

func (c *Client) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	return c.api.V1PaymentIntents.Retrieve(ctx, paymentIntentID, nil)
}

/*
	Private functions
*/

func buildCheckoutSessionParams(request CheckoutSessionRequest, successUrl string, cancelUrl string, now time.Time) *stripe.CheckoutSessionCreateParams {
	lineItem := &stripe.CheckoutSessionCreateLineItemParams{Quantity: stripe.Int64(1)}
	if request.IsUpgrade {
		lineItem.PriceData = &stripe.CheckoutSessionCreateLineItemPriceDataParams{
			Currency:   stripe.String(request.Currency),
			Product:    stripe.String(request.ProductID),
			UnitAmount: stripe.Int64(request.AmountInCents),
		}
	} else {
		lineItem.Price = stripe.String(request.PriceID)
	}
	return &stripe.CheckoutSessionCreateParams{
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:        stripe.String(successUrl),
		CancelURL:         stripe.String(cancelUrl),
		ClientReferenceID: stripe.String(request.UserID),
		CustomerEmail:     stripe.String(request.CustomerEmail),
		LineItems:         []*stripe.CheckoutSessionCreateLineItemParams{lineItem},
		Metadata: map[string]string{
			"transaction_id": request.TransactionID,
			"user_id":        request.UserID,
		},
		PaymentIntentData: &stripe.CheckoutSessionCreatePaymentIntentDataParams{
			Metadata: map[string]string{
				"transaction_id": request.TransactionID,
				"user_id":        request.UserID,
			},
		},
		ExpiresAt: stripe.Int64(now.Add(CheckoutSessionLifetime).Unix()),
	}
}
