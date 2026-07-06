-- +goose Up
-- +goose StatementBegin
ALTER TYPE transaction_status_type ADD VALUE 'expired';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TYPE transaction_status_type_old AS ENUM (
    'pending',
    'completed',
    'failed',
    'refunded'
);

ALTER TABLE transactions
ALTER COLUMN status TYPE transaction_status_type_old
USING status::text::transaction_status_type_old;

DROP TYPE transaction_status_type;

ALTER TYPE transaction_status_type_old
RENAME TO transaction_status_type;
-- +goose StatementEnd