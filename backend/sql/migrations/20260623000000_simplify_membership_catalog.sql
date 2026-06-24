-- +goose Up
-- +goose StatementBegin
ALTER TABLE membership_tiers
    DROP COLUMN is_public,
    DROP COLUMN required_group;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE membership_tiers
    ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN required_group group_type;
-- +goose StatementEnd
