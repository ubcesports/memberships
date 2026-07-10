-- +goose Up
-- +goose StatementBegin
ALTER TABLE membership_tiers
    ADD COLUMN benefits TEXT[] DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE membership_tiers
    DROP COLUMN IF EXISTS benefits;
-- +goose StatementEnd
