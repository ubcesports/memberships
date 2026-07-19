package stripeclient

import (
	"testing"
	"time"
)

func TestBuildCheckoutSessionParamsUsesConfiguredPriceForPurchase(t *testing.T) {
	now := time.Unix(0, 0)
	params := buildCheckoutSessionParams(CheckoutSessionRequest{
		TransactionID: "transaction", UserID: "user", CustomerEmail: "member@example.com",
		PriceID: "price_regular", IsUpgrade: false,
	}, "https://example.com/success", "https://example.com/cancel", now)

	item := params.LineItems[0]
	if item.Price == nil || *item.Price != "price_regular" || item.PriceData != nil {
		t.Fatalf("expected configured Stripe Price, got %#v", item)
	}
	if params.ExpiresAt == nil || *params.ExpiresAt != now.Add(time.Hour).Unix() {
		t.Fatalf("expected one-hour expiration, got %#v", params.ExpiresAt)
	}
}

func TestBuildCheckoutSessionParamsUsesDifferenceForUpgrade(t *testing.T) {
	params := buildCheckoutSessionParams(CheckoutSessionRequest{
		TransactionID: "transaction", UserID: "user", CustomerEmail: "member@example.com",
		PriceID: "price_premium", ProductID: "prod_premium", AmountInCents: 1000,
		Currency: "cad", IsUpgrade: true,
	}, "https://example.com/success", "https://example.com/cancel", time.Unix(0, 0))

	item := params.LineItems[0]
	if item.Price != nil || item.PriceData == nil {
		t.Fatalf("expected inline upgrade PriceData, got %#v", item)
	}
	if item.PriceData.UnitAmount == nil || *item.PriceData.UnitAmount != 1000 {
		t.Fatalf("expected a 1000 minor-unit upgrade, got %#v", item.PriceData.UnitAmount)
	}
	if item.PriceData.Product == nil || *item.PriceData.Product != "prod_premium" {
		t.Fatalf("expected Premium Product, got %#v", item.PriceData.Product)
	}
	if item.PriceData.Currency == nil || *item.PriceData.Currency != "cad" {
		t.Fatalf("expected CAD, got %#v", item.PriceData.Currency)
	}
}
