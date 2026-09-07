-- +goose Up
-- +goose StatementBegin
CREATE TABLE membership_programs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  program_name VARCHAR UNIQUE NOT NULL
);

INSERT INTO membership_programs (program_name)
VALUES ('general')
ON CONFLICT (program_name) DO NOTHING;

ALTER TABLE membership_tiers
    ADD COLUMN program_id UUID;

UPDATE membership_tiers
SET program_id = (
    SELECT id
    FROM membership_programs
    WHERE program_name = 'general'
);

ALTER TABLE membership_tiers
    ALTER COLUMN program_id SET NOT NULL,
    ADD CONSTRAINT membership_tiers_program_id_fkey
        FOREIGN KEY (program_id)
        REFERENCES membership_programs(id);

DROP TRIGGER IF EXISTS trg_one_active_membership_per_user
ON memberships;

DROP FUNCTION IF EXISTS enforce_one_active_membership_per_user();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE membership_tiers
    DROP COLUMN IF EXISTS program_id;

DROP TABLE IF EXISTS membership_programs;
-- +goose StatementEnd
