BEGIN;

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
VALUES
    ('sudipto.islam.mock@example.com', 'MOCK1001', 'member', NOW() - INTERVAL '8 days', NOW(), 'Sudipto Islam', NOW() - INTERVAL '7 days', TRUE, NOW() - INTERVAL '6 days', NULL),
    ('alex.chen.mock@example.com', 'MOCK1002', 'admin', NOW() - INTERVAL '7 days', NOW(), 'Alex Chen', NOW() - INTERVAL '7 days', TRUE, NOW() - INTERVAL '5 days', NULL),
    ('maya.patel.mock@example.com', NULL, 'member', NOW() - INTERVAL '6 days', NOW(), 'Maya Patel', NOW() - INTERVAL '5 days', FALSE, NOW() - INTERVAL '4 days', NULL),
    ('jordan.lee.mock@example.com', 'MOCK1003', 'member', NOW() - INTERVAL '5 days', NOW(), 'Jordan Lee', NULL, TRUE, NULL, NULL),
    ('taylor.morgan.mock@example.com', NULL, 'member', NOW() - INTERVAL '4 days', NOW(), 'Taylor Morgan', NOW() - INTERVAL '3 days', FALSE, NULL, NULL),
    ('priya.shah.mock@example.com', 'MOCK1004', 'admin', NOW() - INTERVAL '3 days', NOW(), 'Priya Shah', NOW() - INTERVAL '3 days', TRUE, NOW() - INTERVAL '2 days', NULL),
    ('noah.williams.mock@example.com', NULL, 'member', NOW() - INTERVAL '2 days', NOW(), 'Noah Williams', NOW() - INTERVAL '1 day', FALSE, NOW() - INTERVAL '1 day', NULL),
    ('emma.garcia.mock@example.com', 'MOCK1005', 'member', NOW() - INTERVAL '1 day', NOW(), 'Emma Garcia', NOW() - INTERVAL '1 day', TRUE, NOW() - INTERVAL '12 hours', NULL)
ON CONFLICT (email) DO UPDATE SET
    student_id = EXCLUDED.student_id,
    role = EXCLUDED.role,
    updated_at = NOW(),
    full_name = EXCLUDED.full_name,
    email_verified_at = EXCLUDED.email_verified_at,
    is_student = EXCLUDED.is_student,
    onboarding_completed_at = EXCLUDED.onboarding_completed_at,
    avatar_url = EXCLUDED.avatar_url;

INSERT INTO user_groups (user_id, "group")
SELECT u.id, mock_groups.group_name::group_type
FROM (
    VALUES
        ('sudipto.islam.mock@example.com', 'member'),
        ('sudipto.islam.mock@example.com', 'competitive_team'),
        ('alex.chen.mock@example.com', 'executive'),
        ('alex.chen.mock@example.com', 'board'),
        ('maya.patel.mock@example.com', 'member'),
        ('jordan.lee.mock@example.com', 'director'),
        ('priya.shah.mock@example.com', 'executive'),
        ('priya.shah.mock@example.com', 'director'),
        ('noah.williams.mock@example.com', 'board'),
        ('emma.garcia.mock@example.com', 'member'),
        ('emma.garcia.mock@example.com', 'competitive_team')
) AS mock_groups(email, group_name)
JOIN users u ON u.email = mock_groups.email
ON CONFLICT (user_id, "group") DO NOTHING;

COMMIT;
