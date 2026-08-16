BEGIN;

-- Mock user with one naturally-expired membership (never cancelled) and one
-- active membership, for exercising the admin panel's "current membership" +
-- "membership history" views together.
INSERT INTO users (
    email,
    student_id,
    role,
    created_at,
    updated_at,
    full_name,
    email_verified_at,
    is_student,
    onboarding_completed_at,
    avatar_url
)
VALUES (
    'jamie.rivera.mock@example.com',
    'MOCK1006',
    'member',
    NOW() - INTERVAL '400 days',
    NOW(),
    'Jamie Rivera',
    NOW() - INTERVAL '400 days',
    TRUE,
    NOW() - INTERVAL '399 days',
    NULL
)
ON CONFLICT (email) DO UPDATE SET
    student_id = EXCLUDED.student_id,
    role = EXCLUDED.role,
    updated_at = NOW(),
    full_name = EXCLUDED.full_name,
    email_verified_at = EXCLUDED.email_verified_at,
    is_student = EXCLUDED.is_student,
    onboarding_completed_at = EXCLUDED.onboarding_completed_at;

-- Keep this seed rerunnable without duplicating rows.
DELETE FROM transactions
WHERE user_id = (SELECT id FROM users WHERE email = 'jamie.rivera.mock@example.com')
    AND stripe_payment_intent_id LIKE 'seed_mock_pi_jamie_%';

DELETE FROM memberships
WHERE user_id = (SELECT id FROM users WHERE email = 'jamie.rivera.mock@example.com');

WITH expired_membership AS (
    INSERT INTO memberships (user_id, tier_id, started_at, expires_at, cancelled_at)
    SELECT u.id, mt.id, NOW() - INTERVAL '400 days', NOW() - INTERVAL '35 days', NULL
    FROM users u
    JOIN membership_tiers mt ON mt.slug = 'basic'
    WHERE u.email = 'jamie.rivera.mock@example.com'
    RETURNING id, user_id, tier_id
),
active_membership AS (
    INSERT INTO memberships (user_id, tier_id, started_at, expires_at, cancelled_at)
    SELECT u.id, mt.id, NOW() - INTERVAL '30 days', NOW() + INTERVAL '300 days', NULL
    FROM users u
    JOIN membership_tiers mt ON mt.slug = 'lounge'
    WHERE u.email = 'jamie.rivera.mock@example.com'
    RETURNING id, user_id, tier_id
)
INSERT INTO transactions (
    user_id,
    membership_id,
    tier_id,
    stripe_payment_intent_id,
    status,
    group_at_purchase,
    student_at_purchase,
    amount_paid_cents,
    purchase_type
)
SELECT user_id, id, tier_id, 'seed_mock_pi_jamie_expired', 'completed'::transaction_status_type, 'member'::group_type, TRUE, 2500, 'new'::purchase_type
FROM expired_membership
UNION ALL
SELECT user_id, id, tier_id, 'seed_mock_pi_jamie_active', 'completed'::transaction_status_type, 'member'::group_type, TRUE, 4500, 'new'::purchase_type
FROM active_membership;

COMMIT;
