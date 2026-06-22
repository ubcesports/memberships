-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN avatar_url TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN IF EXISTS avatar_url;
-- +goose StatementEnd
