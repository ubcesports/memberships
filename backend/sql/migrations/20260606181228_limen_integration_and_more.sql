-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS verification_tokens;
DROP TABLE IF EXISTS transactions;
DROP INDEX IF EXISTS memberships_one_active_per_user;
DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS promo_codes;

DROP TYPE IF EXISTS verification_token_type;
DROP TYPE IF EXISTS membership_status_type;

ALTER TABLE users ALTER COLUMN role DROP DEFAULT;
UPDATE users SET role = 'member' WHERE role NOT IN ('member', 'admin');
ALTER TABLE users ALTER COLUMN role TYPE VARCHAR;
DROP TYPE IF EXISTS role_type;
CREATE TYPE role_type AS ENUM ('member', 'admin');
ALTER TABLE users ALTER COLUMN role TYPE role_type USING role::role_type;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'member';
ALTER TABLE users ALTER COLUMN role SET NOT NULL;

CREATE TYPE group_type AS ENUM (
    'member',
    'competitive_team',
    'executive',
    'director',
    'board'
);

ALTER TABLE users
    DROP COLUMN is_verified,
    DROP COLUMN first_name,
    DROP COLUMN last_name,
    ADD COLUMN full_name VARCHAR NOT NULL,
    ADD COLUMN email_verified_at TIMESTAMPTZ,
    ADD COLUMN password VARCHAR,
    ADD COLUMN is_student BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN onboarding_completed_at TIMESTAMPTZ;

ALTER TABLE users RENAME COLUMN ubc_student_id TO student_id;

DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    last_access TIMESTAMPTZ NOT NULL,
    metadata TEXT
);

DROP TABLE IF EXISTS verification_tokens;
CREATE TABLE verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject VARCHAR NOT NULL,
    value VARCHAR NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    provider VARCHAR NOT NULL,
    provider_account_id VARCHAR,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    access_token_expires_at TIMESTAMPTZ,
    scope VARCHAR NOT NULL,
    id_token TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (provider, provider_account_id)
);

CREATE TABLE user_groups (
    user_id UUID NOT NULL REFERENCES users(id),
    "group" group_type NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, "group")
);

ALTER TABLE membership_tiers
    DROP COLUMN price_student,
    DROP COLUMN price_non_student,
    ADD COLUMN stripe_product_id VARCHAR UNIQUE,
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

CREATE TABLE membership_tier_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tier_id UUID NOT NULL REFERENCES membership_tiers(id),
    "group" group_type NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    stripe_price_id VARCHAR UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tier_id, "group")
);

DROP TABLE IF EXISTS transactions;
DROP INDEX IF EXISTS memberships_one_active_per_user;
DROP TABLE IF EXISTS memberships;
DROP TYPE IF EXISTS membership_status_type;

CREATE TABLE memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    tier_id UUID NOT NULL REFERENCES membership_tiers(id),
    transaction_id UUID,
    group_at_purchase group_type NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    membership_id UUID NOT NULL REFERENCES memberships(id),
    stripe_payment_intent_id VARCHAR UNIQUE,
    price_amount NUMERIC(10,2) NOT NULL,
    status transaction_status_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE memberships
    ADD CONSTRAINT memberships_transaction_id_fkey
    FOREIGN KEY (transaction_id) REFERENCES transactions(id);

CREATE UNIQUE INDEX memberships_one_active_per_user
    ON memberships (user_id)
    WHERE cancelled_at IS NULL;

DROP TABLE IF EXISTS promo_codes;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS memberships DROP CONSTRAINT IF EXISTS memberships_transaction_id_fkey;

DROP TABLE IF EXISTS transactions;
DROP INDEX IF EXISTS memberships_one_active_per_user;
DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS membership_tier_prices;
DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS verifications;
DROP TABLE IF EXISTS sessions;

CREATE TYPE verification_token_type AS ENUM ('email_verification', 'magic_link');
CREATE TYPE membership_status_type AS ENUM ('active', 'expired', 'cancelled');

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ DEFAULT NOW(),
    ip_address VARCHAR,
    user_agent TEXT
);

CREATE TABLE verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    token VARCHAR NOT NULL UNIQUE,
    type verification_token_type NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE membership_tiers
    DROP COLUMN IF EXISTS stripe_product_id,
    DROP COLUMN IF EXISTS is_active,
    ADD COLUMN price_student NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN price_non_student NUMERIC(10,2) NOT NULL DEFAULT 0;

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

ALTER TABLE users
    DROP COLUMN IF EXISTS email_verified_at,
    DROP COLUMN IF EXISTS password,
    DROP COLUMN IF EXISTS is_student,
    DROP COLUMN IF EXISTS onboarding_completed_at,
    DROP COLUMN IF EXISTS full_name,
    ADD COLUMN first_name VARCHAR NOT NULL DEFAULT '',
    ADD COLUMN last_name VARCHAR NOT NULL DEFAULT '',
    ADD COLUMN is_verified BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE users RENAME COLUMN student_id TO ubc_student_id;

ALTER TABLE users ALTER COLUMN role DROP DEFAULT;
ALTER TABLE users ALTER COLUMN role TYPE VARCHAR;
DROP TYPE IF EXISTS role_type;
CREATE TYPE role_type AS ENUM ('member', 'exec', 'competitive', 'admin');

CREATE TABLE promo_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR NOT NULL UNIQUE,
    discount_percentage SMALLINT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    role_override role_type DEFAULT 'member',
    CONSTRAINT chk_discount_range CHECK (discount_percentage >= 0 AND discount_percentage <= 100)
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

ALTER TABLE users ALTER COLUMN role TYPE role_type USING role::role_type;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'member';

DROP TYPE IF EXISTS group_type;
-- +goose StatementEnd