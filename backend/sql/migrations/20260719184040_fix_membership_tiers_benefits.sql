-- +goose Up
-- +goose StatementBegin
ALTER TABLE memberships
    DROP COLUMN IF EXISTS benefits;

ALTER TABLE membership_tiers
    ADD COLUMN IF NOT EXISTS benefits TEXT[] NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE membership_tiers
    DROP COLUMN IF EXISTS benefits;

ALTER TABLE memberships
    ADD COLUMN IF NOT EXISTS benefits TEXT[] NOT NULL DEFAULT '{}';
-- +goose StatementEnd