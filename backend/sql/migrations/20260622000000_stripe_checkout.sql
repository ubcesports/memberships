-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE membership_tiers
    ADD COLUMN slug VARCHAR NOT NULL,
    ADD CONSTRAINT membership_tiers_slug_key UNIQUE (slug);

ALTER TABLE membership_tier_prices
    DROP CONSTRAINT IF EXISTS membership_tier_prices_tier_id_group_key,
    ADD COLUMN is_student BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE membership_tier_prices
    DROP COLUMN price,
    ALTER COLUMN stripe_price_id SET NOT NULL,
    ADD CONSTRAINT membership_tier_prices_tier_group_student_key
        UNIQUE (tier_id, "group", is_student);

DROP INDEX IF EXISTS memberships_one_active_per_user;
ALTER TABLE memberships
    DROP CONSTRAINT IF EXISTS memberships_transaction_id_fkey,
    DROP COLUMN IF EXISTS transaction_id,
    ADD COLUMN is_student_at_purchase BOOLEAN NOT NULL DEFAULT FALSE,
    ADD CONSTRAINT memberships_valid_period CHECK (started_at < expires_at),
    ADD CONSTRAINT memberships_no_overlapping_active_periods
        EXCLUDE USING gist (
            user_id WITH =,
            tstzrange(started_at, expires_at, '[)') WITH &&
        ) WHERE (cancelled_at IS NULL);

ALTER TABLE transactions
    ALTER COLUMN membership_id DROP NOT NULL;

ALTER TABLE transactions
    RENAME COLUMN price_amount TO amount_minor;

ALTER TABLE transactions
    ALTER COLUMN amount_minor TYPE BIGINT USING ROUND(amount_minor * 100)::BIGINT,
    ADD COLUMN tier_id UUID REFERENCES membership_tiers(id),
    ADD COLUMN group_at_purchase group_type,
    ADD COLUMN is_student_at_purchase BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN stripe_checkout_session_id VARCHAR UNIQUE,
    ADD COLUMN stripe_charge_id VARCHAR UNIQUE,
    ADD COLUMN stripe_price_id VARCHAR,
    ADD COLUMN currency VARCHAR(3),
    ADD CONSTRAINT transactions_amount_nonnegative CHECK (amount_minor >= 0);

CREATE UNIQUE INDEX transactions_one_pending_per_user
    ON transactions (user_id)
    WHERE status = 'pending';

CREATE INDEX transactions_payment_intent_idx
    ON transactions (stripe_payment_intent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS transactions_payment_intent_idx;
DROP INDEX IF EXISTS transactions_one_pending_per_user;

ALTER TABLE transactions
    DROP CONSTRAINT IF EXISTS transactions_amount_nonnegative,
    DROP COLUMN IF EXISTS currency,
    DROP COLUMN IF EXISTS stripe_price_id,
    DROP COLUMN IF EXISTS stripe_charge_id,
    DROP COLUMN IF EXISTS stripe_checkout_session_id,
    DROP COLUMN IF EXISTS is_student_at_purchase,
    DROP COLUMN IF EXISTS group_at_purchase,
    DROP COLUMN IF EXISTS tier_id;

ALTER TABLE transactions
    RENAME COLUMN amount_minor TO price_amount;

ALTER TABLE transactions
    ALTER COLUMN price_amount TYPE NUMERIC(10,2) USING (price_amount::NUMERIC / 100),
    ALTER COLUMN membership_id SET NOT NULL;

ALTER TABLE memberships
    DROP CONSTRAINT IF EXISTS memberships_no_overlapping_active_periods,
    DROP CONSTRAINT IF EXISTS memberships_valid_period,
    DROP COLUMN IF EXISTS is_student_at_purchase,
    ADD COLUMN transaction_id UUID REFERENCES transactions(id);

CREATE UNIQUE INDEX memberships_one_active_per_user
    ON memberships (user_id)
    WHERE cancelled_at IS NULL;

ALTER TABLE membership_tier_prices
    DROP CONSTRAINT IF EXISTS membership_tier_prices_tier_group_student_key,
    DROP COLUMN IF EXISTS is_student,
    ADD COLUMN price NUMERIC(10,2) NOT NULL DEFAULT 0,
    ALTER COLUMN stripe_price_id DROP NOT NULL,
    ADD CONSTRAINT membership_tier_prices_tier_id_group_key UNIQUE (tier_id, "group");

ALTER TABLE membership_tiers
    DROP CONSTRAINT IF EXISTS membership_tiers_slug_key,
    DROP COLUMN IF EXISTS slug;
-- +goose StatementEnd
