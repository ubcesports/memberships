-- +goose Up
-- +goose StatementBegin
CREATE TYPE tier_code_type AS ENUM (
    'regular',
    'premium',
    'cab',
    'day'
);

CREATE TYPE student_status_type AS ENUM (
    'student',
    'non_student'
);

ALTER TABLE membership_tiers
    ADD COLUMN code tier_code_type NOT NULL UNIQUE;

ALTER TABLE membership_tier_prices
    ADD COLUMN student_status student_status_type NOT NULL;

ALTER TABLE membership_tier_prices
    DROP CONSTRAINT membership_tier_prices_tier_id_group_key;

ALTER TABLE membership_tier_prices
    ADD CONSTRAINT membership_tier_prices_unique_tier_group_student
    UNIQUE (tier_id, "group", student_status);

ALTER TABLE memberships
    ADD COLUMN student_status_at_purchase student_status_type NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE memberships
    DROP COLUMN IF EXISTS student_status_at_purchase;

ALTER TABLE membership_tier_prices
    DROP CONSTRAINT IF EXISTS membership_tier_prices_unique_tier_group_student;

ALTER TABLE membership_tier_prices
    DROP COLUMN IF EXISTS student_status;

ALTER TABLE membership_tier_prices
    ADD CONSTRAINT membership_tier_prices_tier_id_group_key
    UNIQUE (tier_id, "group");

ALTER TABLE membership_tiers
    DROP COLUMN IF EXISTS code;

DROP TYPE IF EXISTS student_status_type;
DROP TYPE IF EXISTS tier_code_type;
-- +goose StatementEnd
