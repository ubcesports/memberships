-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions
DROP CONSTRAINT IF EXISTS transactions_membership_id_fkey;

ALTER TABLE transactions
DROP COLUMN IF EXISTS membership_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
ADD COLUMN IF NOT EXISTS membership_id uuid;

ALTER TABLE transactions
ADD CONSTRAINT transactions_membership_id_fkey
FOREIGN KEY (membership_id)
REFERENCES memberships(id);
-- +goose StatementEnd