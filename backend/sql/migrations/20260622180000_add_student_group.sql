-- +goose Up
-- +goose StatementBegin
ALTER TYPE group_type ADD VALUE IF NOT EXISTS 'student';
-- +goose StatementEnd

-- +goose Down
SELECT 1;
