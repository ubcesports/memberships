-- +goose Up
-- +goose StatementBegin
ALTER TYPE transaction_status_type ADD VALUE IF NOT EXISTS 'expired';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE transactions SET status = 'failed' WHERE status = 'expired';

DROP INDEX transactions_one_pending_per_user;

ALTER TABLE transactions
    ALTER COLUMN status TYPE VARCHAR USING status::text;

DROP TYPE transaction_status_type;

CREATE TYPE transaction_status_type AS ENUM (
    'pending',
    'completed',
    'failed',
    'refunded'
);

ALTER TABLE transactions
    ALTER COLUMN status TYPE transaction_status_type
    USING status::transaction_status_type;

CREATE UNIQUE INDEX transactions_one_pending_per_user
    ON transactions (user_id)
    WHERE status = 'pending';
-- +goose StatementEnd
