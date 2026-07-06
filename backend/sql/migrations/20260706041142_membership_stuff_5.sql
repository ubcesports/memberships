-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions
    ADD COLUMN stripe_checkout_session_id VARCHAR UNIQUE;

CREATE UNIQUE INDEX one_pending_transaction_per_user
ON transactions (user_id)
WHERE status = 'pending';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX one_pending_transaction_per_user;

ALTER TABLE transactions
    DROP COLUMN stripe_checkout_session_id;
-- +goose StatementEnd
