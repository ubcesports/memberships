-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS student_at_purchase BOOLEAN,
    DROP COLUMN IF EXISTS price_amount,
    ADD COLUMN IF NOT EXISTS amount_paid_cents BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
    DROP COLUMN IF EXISTS student_at_purchase,
    ADD COLUMN IF NOT EXISTS price_amount NUMERIC(10, 2),
    DROP COLUMN IF EXISTS amount_paid_cents;
-- +goose StatementEnd