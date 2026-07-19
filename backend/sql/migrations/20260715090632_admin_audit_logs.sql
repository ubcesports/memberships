-- +goose Up
-- +goose StatementBegin
CREATE TYPE admin_audit_outcome_type AS ENUM (
    'success',
    'failed',
    'denied'
);

CREATE TABLE admin_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    actor_user_id UUID NOT NULL REFERENCES users(id),
    action VARCHAR NOT NULL,
    target_user_id UUID REFERENCES users(id),
    outcome admin_audit_outcome_type NOT NULL,
    request_id VARCHAR NOT NULL,
    description TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE admin_audit_logs;
DROP TYPE admin_audit_outcome_type;
-- +goose StatementEnd
