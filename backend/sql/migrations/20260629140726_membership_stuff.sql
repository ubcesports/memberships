-- +goose Up
-- +goose StatementBegin
ALTER TABLE memberships
    DROP COLUMN IF EXISTS group_at_purchase,
    DROP COLUMN IF EXISTS transaction_id;

ALTER TABLE membership_tiers
    ADD COLUMN IF NOT EXISTS slug VARCHAR,
    ADD COLUMN IF NOT EXISTS "group" group_type;

ALTER TABLE membership_tier_prices
    DROP COLUMN IF EXISTS price,
    ADD COLUMN IF NOT EXISTS is_student_required BOOLEAN,
    DROP COLUMN IF EXISTS "group";

ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS group_at_purchase group_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE memberships
    ADD COLUMN IF NOT EXISTS group_at_purchase group_type,
    ADD COLUMN IF NOT EXISTS transaction_id UUID;

ALTER TABLE membership_tiers
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS "group";

ALTER TABLE membership_tier_prices
    ADD COLUMN IF NOT EXISTS price NUMERIC(10, 2),
    DROP COLUMN IF EXISTS is_student_required,
    ADD COLUMN IF NOT EXISTS "group" group_type;

ALTER TABLE transactions
    DROP COLUMN IF EXISTS group_at_purchase;
-- +goose StatementEnd