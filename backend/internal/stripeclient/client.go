package stripeclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/stripe/stripe-go/v85"
)

type Gateway interface {
	GetPrice(context.Context, string) (*stripe.Price, error)
	CreateCheckoutSession(context.Context, CheckoutSessionRequest) (*stripe.CheckoutSession, error)
	GetCheckoutSession(context.Context, string) (*stripe.CheckoutSession, error)
	ExpireCheckoutSession(context.Context, string) (*stripe.CheckoutSession, error)
	GetPaymentIntent(context.Context, string) (*stripe.PaymentIntent, error)
}

type CheckoutSessionRequest struct {
	TransactionID string
	UserID        string
	CustomerEmail string
	PriceID       string
}

type Client struct {
	api        *stripe.Client
	successURL string
	cancelURL  string
}

func NewClient() (*Client, error) {
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	successURL := os.Getenv("STRIPE_CHECKOUT_SUCCESS_URL")
	cancelURL := os.Getenv("STRIPE_CHECKOUT_CANCEL_URL")
	if secretKey == "" || successURL == "" || cancelURL == "" {
		return nil, fmt.Errorf("STRIPE_SECRET_KEY, STRIPE_CHECKOUT_SUCCESS_URL, and STRIPE_CHECKOUT_CANCEL_URL are required")
	}

	return &Client{
		api:        stripe.NewClient(secretKey),
		successURL: successURL,
		cancelURL:  cancelURL,
	}, nil
}

func (c *Client) GetPrice(ctx context.Context, priceID string) (*stripe.Price, error) {
	return c.api.V1Prices.Retrieve(ctx, priceID, nil)
}

func (c *Client) CreateCheckoutSession(ctx context.Context, request CheckoutSessionRequest) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionCreateParams{
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:        stripe.String(c.successURL),
		CancelURL:         stripe.String(c.cancelURL),
		ClientReferenceID: stripe.String(request.UserID),
		CustomerEmail:     stripe.String(request.CustomerEmail),
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{Price: stripe.String(request.PriceID), Quantity: stripe.Int64(1)},
		},
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
		ExpiresAt: stripe.Int64(time.Now().Add(30 * time.Minute).Unix()),
	}
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
