-- +goose Up
-- +goose StatementBegin
CREATE TYPE payment_method_type AS ENUM (
    'stripe',
    'cash',
    'etransfer'
);

ALTER TABLE transactions
    ADD COLUMN payment_method payment_method_type NOT NULL DEFAULT 'stripe';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
    DROP COLUMN IF EXISTS payment_method;

DROP TYPE IF EXISTS payment_method_type;
-- +goose StatementEnd
