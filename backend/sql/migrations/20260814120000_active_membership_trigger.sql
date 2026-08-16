-- +goose Up
-- +goose StatementBegin

-- The old partial unique index (WHERE cancelled_at IS NULL) blocked a user
-- from ever having a naturally-expired, never-cancelled membership sit
-- alongside a newer active one, since both have cancelled_at IS NULL.
-- CancelActiveMembershipsByUserId only cancels memberships that are
-- currently within their active date window, so a membership that already
-- expired on its own is left with cancelled_at NULL forever. Replace the
-- index with a trigger that only rejects a second *currently active*
-- (cancelled_at IS NULL AND expires_at > NOW()) membership per user.
DROP INDEX IF EXISTS memberships_one_active_per_user;

CREATE OR REPLACE FUNCTION enforce_one_active_membership_per_user() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.cancelled_at IS NULL AND NEW.expires_at > NOW() THEN
        IF EXISTS (
            SELECT 1 FROM memberships
            WHERE user_id = NEW.user_id
                AND id <> NEW.id
                AND cancelled_at IS NULL
                AND expires_at > NOW()
        ) THEN
            RAISE EXCEPTION 'user % already has an active membership', NEW.user_id
                USING ERRCODE = 'unique_violation';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_one_active_membership_per_user
BEFORE INSERT OR UPDATE ON memberships
FOR EACH ROW EXECUTE FUNCTION enforce_one_active_membership_per_user();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trg_one_active_membership_per_user ON memberships;
DROP FUNCTION IF EXISTS enforce_one_active_membership_per_user();

CREATE UNIQUE INDEX memberships_one_active_per_user
    ON memberships (user_id)
    WHERE cancelled_at IS NULL;

-- +goose StatementEnd
