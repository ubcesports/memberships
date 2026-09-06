-- +goose Up
-- +goose StatementBegin
CREATE TYPE membership_expiration_type AS ENUM (
    'semester',
    'year',
    'day'
);

ALTER TABLE membership_tiers
    ADD COLUMN expiration_type membership_expiration_type NOT NULL DEFAULT 'semester';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE membership_tiers
    DROP COLUMN IF EXISTS expiration_type;

DROP TYPE IF EXISTS membership_expiration_type;
-- +goose StatementEnd
