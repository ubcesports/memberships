-- +goose Up
-- +goose StatementBegin
CREATE TYPE purchase_type AS ENUM (
    'new',
    'upgrade'
);

ALTER TYPE transaction_status_type ADD VALUE IF NOT EXISTS 'expired';

ALTER TABLE memberships
    DROP COLUMN IF EXISTS group_at_purchase,
    ADD COLUMN benefits TEXT[] DEFAULT '{}',
    DROP COLUMN IF EXISTS transaction_id;

ALTER TABLE membership_tiers
    ADD COLUMN IF NOT EXISTS slug VARCHAR,
    ADD COLUMN IF NOT EXISTS "group" group_type;

ALTER TABLE membership_tier_prices
    DROP COLUMN IF EXISTS price,
    ADD COLUMN IF NOT EXISTS is_student_required BOOLEAN,
    DROP COLUMN IF EXISTS "group";

ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS group_at_purchase group_type,
    ADD COLUMN IF NOT EXISTS student_at_purchase BOOLEAN,
    DROP COLUMN IF EXISTS price_amount,
    ADD COLUMN IF NOT EXISTS amount_paid_cents BIGINT,
    ADD COLUMN IF NOT EXISTS purchase_type purchase_type,
    ADD COLUMN IF NOT EXISTS stripe_checkout_session_id VARCHAR UNIQUE,
    ADD COLUMN IF NOT EXISTS tier_id UUID REFERENCES membership_tiers(id),
    ALTER COLUMN membership_id DROP NOT NULL;

CREATE UNIQUE INDEX one_pending_transaction_per_user
ON transactions (user_id)
WHERE status = 'pending';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS one_pending_transaction_per_user;

ALTER TABLE transactions
    DROP COLUMN IF EXISTS stripe_checkout_session_id,
    DROP COLUMN IF EXISTS tier_id,
    DROP COLUMN IF EXISTS purchase_type,
    DROP COLUMN IF EXISTS group_at_purchase,
    DROP COLUMN IF EXISTS student_at_purchase,
    ADD COLUMN IF NOT EXISTS price_amount NUMERIC(10, 2),
    DROP COLUMN IF EXISTS amount_paid_cents,
    ALTER COLUMN membership_id SET NOT NULL;

DROP TYPE IF EXISTS purchase_type;

ALTER TABLE membership_tiers
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS benefits,
    DROP COLUMN IF EXISTS "group";

ALTER TABLE membership_tier_prices
    ADD COLUMN IF NOT EXISTS price NUMERIC(10, 2),
    DROP COLUMN IF EXISTS is_student_required,
    ADD COLUMN IF NOT EXISTS "group" group_type;

ALTER TABLE memberships
    ADD COLUMN IF NOT EXISTS group_at_purchase group_type,
    ADD COLUMN IF NOT EXISTS transaction_id UUID;

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
