-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions
    ADD COLUMN tier_id UUID REFERENCES membership_tiers(id),
    ALTER COLUMN membership_id DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
    DROP COLUMN IF EXISTS tier_id,
    ALTER COLUMN membership_id SET NOT NULL;
-- +goose StatementEnd
