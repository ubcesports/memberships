-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE role_type AS ENUM (
    'member',
    'exec',
    'competitive',
    'admin'
);

CREATE TYPE verification_token_type AS ENUM (
    'email_verification',
    'magic_link'
);

CREATE TYPE membership_status_type AS ENUM (
    'active',
    'expired',
    'cancelled'
);

CREATE TYPE transaction_status_type AS ENUM (
    'pending',
    'completed',
    'failed',
    'refunded'
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR NOT NULL UNIQUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    ubc_student_id VARCHAR UNIQUE,
    role role_type NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    token VARCHAR NOT NULL UNIQUE,
    type verification_token_type NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ,
    ip_address VARCHAR,
    user_agent TEXT
);

CREATE TABLE membership_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR NOT NULL,
    description TEXT,
    price_student NUMERIC(10,2) NOT NULL,
    price_non_student NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    tier_id UUID NOT NULL REFERENCES membership_tiers(id),
    status membership_status_type NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX memberships_one_active_per_user
ON memberships (user_id)
WHERE status = 'active';

CREATE TABLE promo_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR NOT NULL UNIQUE,
    discount_percentage SMALLINT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    role_override role_type DEFAULT 'member',
    CONSTRAINT chk_discount_range
        CHECK (
            discount_percentage >= 0
            AND discount_percentage <= 100
        )
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    membership_id UUID NOT NULL REFERENCES memberships(id),
    promo_code_id UUID REFERENCES promo_codes(id),
    price_amount NUMERIC(10,2) NOT NULL,
    status transaction_status_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS promo_codes;

DROP INDEX IF EXISTS memberships_one_active_per_user;

DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS membership_tiers;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS verification_tokens;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS transaction_status_type;
DROP TYPE IF EXISTS membership_status_type;
DROP TYPE IF EXISTS verification_token_type;
DROP TYPE IF EXISTS role_type;
-- +goose StatementEnd