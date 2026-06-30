-- +goose Up
-- +goose StatementBegin
CREATE TYPE purchase_type AS ENUM (
    'new',
    'upgrade'
);

ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS purchase_type purchase_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
    DROP COLUMN IF EXISTS purchase_type;

DROP TYPE IF EXISTS purchase_type;
-- +goose StatementEnd