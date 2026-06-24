-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION assign_default_member_group()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO user_groups (user_id, "group")
    VALUES (NEW.id, 'member')
    ON CONFLICT (user_id, "group") DO NOTHING;

    RETURN NEW;
END;
$$;

CREATE TRIGGER users_assign_default_member_group
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION assign_default_member_group();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS users_assign_default_member_group ON users;
DROP FUNCTION IF EXISTS assign_default_member_group();
-- +goose StatementEnd
