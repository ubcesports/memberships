-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_student;

ALTER TABLE membership_tiers
    ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN required_group group_type;

ALTER TABLE membership_tier_prices
    DROP CONSTRAINT IF EXISTS membership_tier_prices_tier_group_student_key;

ALTER TABLE membership_tier_prices
    DROP COLUMN is_student,
    ADD CONSTRAINT membership_tier_prices_tier_group_key UNIQUE (tier_id, "group");

ALTER TABLE memberships DROP COLUMN is_student_at_purchase;

CREATE TYPE transaction_kind_type AS ENUM ('purchase', 'upgrade');

ALTER TABLE transactions
    DROP COLUMN is_student_at_purchase,
    ADD COLUMN kind transaction_kind_type NOT NULL DEFAULT 'purchase',
    ADD COLUMN credit_amount_minor BIGINT NOT NULL DEFAULT 0,
    ADD CONSTRAINT transactions_credit_nonnegative CHECK (credit_amount_minor >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
    DROP CONSTRAINT IF EXISTS transactions_credit_nonnegative,
    DROP COLUMN IF EXISTS credit_amount_minor,
    DROP COLUMN IF EXISTS kind,
    ADD COLUMN is_student_at_purchase BOOLEAN NOT NULL DEFAULT FALSE;

DROP TYPE transaction_kind_type;

ALTER TABLE memberships
    ADD COLUMN is_student_at_purchase BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE membership_tier_prices
    DROP CONSTRAINT IF EXISTS membership_tier_prices_tier_group_key,
    ADD COLUMN is_student BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE membership_tier_prices
    ADD CONSTRAINT membership_tier_prices_tier_group_student_key
        UNIQUE (tier_id, "group", is_student);

ALTER TABLE membership_tiers
    DROP COLUMN IF EXISTS required_group,
    DROP COLUMN IF EXISTS is_public;

ALTER TABLE users ADD COLUMN is_student BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd
