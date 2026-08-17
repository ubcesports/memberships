-- +goose Up
-- +goose StatementBegin
ALTER TABLE membership_tier_prices
    ADD COLUMN price_in_cents BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE membership_tier_prices
    DROP COLUMN IF EXISTS price_in_cents;
-- +goose StatementEnd
