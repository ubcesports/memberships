-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TYPE group_type ADD VALUE IF NOT EXISTS 'student';
ALTER TYPE transaction_status_type ADD VALUE IF NOT EXISTS 'expired';
CREATE TYPE transaction_kind_type AS ENUM ('purchase', 'upgrade', 'replacement');

ALTER TABLE users
    DROP COLUMN is_student;

ALTER TABLE membership_tiers
    ADD COLUMN slug VARCHAR NOT NULL,
    ADD CONSTRAINT membership_tiers_slug_key UNIQUE (slug);

ALTER TABLE membership_tier_prices
    DROP COLUMN price,
    ALTER COLUMN stripe_price_id SET NOT NULL;

DROP INDEX IF EXISTS memberships_one_active_per_user;

ALTER TABLE memberships
    DROP CONSTRAINT IF EXISTS memberships_transaction_id_fkey,
    DROP COLUMN transaction_id,
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
    ADD COLUMN stripe_checkout_session_id VARCHAR UNIQUE,
    ADD COLUMN stripe_charge_id VARCHAR UNIQUE,
    ADD COLUMN stripe_price_id VARCHAR,
    ADD COLUMN currency VARCHAR(3),
    ADD COLUMN kind transaction_kind_type NOT NULL DEFAULT 'purchase',
    ADD COLUMN credit_amount_minor BIGINT NOT NULL DEFAULT 0,
    ADD CONSTRAINT transactions_amount_nonnegative CHECK (amount_minor >= 0),
    ADD CONSTRAINT transactions_credit_nonnegative CHECK (credit_amount_minor >= 0);

CREATE UNIQUE INDEX transactions_one_pending_per_user
    ON transactions (user_id)
    WHERE status = 'pending';

CREATE INDEX transactions_payment_intent_idx
    ON transactions (stripe_payment_intent_id);

CREATE FUNCTION assign_default_member_group()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO user_groups (user_id, "group")
    VALUES (NEW.id, 'member')
    ON CONFLICT (user_id, "group") DO NOTHING;

    RETURN NEW;
END;
$$;

CREATE TRIGGER users_assign_default_member_group
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION assign_default_member_group();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS users_assign_default_member_group ON users;
DROP FUNCTION IF EXISTS assign_default_member_group();

DROP INDEX IF EXISTS transactions_payment_intent_idx;
DROP INDEX IF EXISTS transactions_one_pending_per_user;

ALTER TABLE transactions
    DROP CONSTRAINT IF EXISTS transactions_credit_nonnegative,
    DROP CONSTRAINT IF EXISTS transactions_amount_nonnegative,
    DROP COLUMN credit_amount_minor,
    DROP COLUMN kind,
    DROP COLUMN currency,
    DROP COLUMN stripe_price_id,
    DROP COLUMN stripe_charge_id,
    DROP COLUMN stripe_checkout_session_id,
    DROP COLUMN group_at_purchase,
    DROP COLUMN tier_id;

DROP TYPE transaction_kind_type;

ALTER TABLE transactions
    RENAME COLUMN amount_minor TO price_amount;

ALTER TABLE transactions
    ALTER COLUMN price_amount TYPE NUMERIC(10,2) USING (price_amount::NUMERIC / 100),
    ALTER COLUMN membership_id SET NOT NULL;

ALTER TABLE memberships
    DROP CONSTRAINT IF EXISTS memberships_no_overlapping_active_periods,
    DROP CONSTRAINT IF EXISTS memberships_valid_period,
    ADD COLUMN transaction_id UUID REFERENCES transactions(id);

CREATE UNIQUE INDEX memberships_one_active_per_user
    ON memberships (user_id)
    WHERE cancelled_at IS NULL;

ALTER TABLE membership_tier_prices
    ADD COLUMN price NUMERIC(10,2) NOT NULL DEFAULT 0,
    ALTER COLUMN stripe_price_id DROP NOT NULL;

ALTER TABLE membership_tiers
    DROP CONSTRAINT IF EXISTS membership_tiers_slug_key,
    DROP COLUMN slug;

ALTER TABLE users
    ADD COLUMN is_student BOOLEAN NOT NULL DEFAULT FALSE;

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

ALTER TABLE user_groups
    ALTER COLUMN "group" TYPE VARCHAR USING "group"::text;

ALTER TABLE memberships
    ALTER COLUMN group_at_purchase TYPE VARCHAR USING group_at_purchase::text;

ALTER TABLE membership_tier_prices
    ALTER COLUMN "group" TYPE VARCHAR USING "group"::text;

DROP TYPE group_type;

CREATE TYPE group_type AS ENUM (
    'member',
    'competitive_team',
    'executive',
    'director',
    'board'
);

ALTER TABLE user_groups
    ALTER COLUMN "group" TYPE group_type USING "group"::group_type;

ALTER TABLE memberships
    ALTER COLUMN group_at_purchase TYPE group_type USING group_at_purchase::group_type;

ALTER TABLE membership_tier_prices
    ALTER COLUMN "group" TYPE group_type USING "group"::group_type;
-- +goose StatementEnd
