BEGIN;

-- These fixed-ID users satisfy the audit log foreign keys and keep this seed
-- independent from other seed files.
INSERT INTO users (
    id,
    email,
    role,
    created_at,
    updated_at,
    full_name,
    is_student,
    avatar_url
)
VALUES
    (
        'a1000000-0000-4000-8000-000000000001',
        'audit-admin-alex.mock@example.com',
        'admin',
        NOW() - INTERVAL '30 days',
        NOW(),
        'Alex Chen',
        TRUE,
        'https://api.dicebear.com/9.x/initials/svg?seed=Alex%20Chen'
    ),
    (
        'a2000000-0000-4000-8000-000000000002',
        'audit-admin-priya.mock@example.com',
        'admin',
        NOW() - INTERVAL '25 days',
        NOW(),
        'Priya Shah',
        TRUE,
        'https://api.dicebear.com/9.x/initials/svg?seed=Priya%20Shah'
    ),
    (
        'b1000000-0000-4000-8000-000000000001',
        'audit-member-maya.mock@example.com',
        'member',
        NOW() - INTERVAL '20 days',
        NOW(),
        'Maya Patel',
        FALSE,
        'https://api.dicebear.com/9.x/initials/svg?seed=Maya%20Patel'
    ),
    (
        'b2000000-0000-4000-8000-000000000002',
        'audit-member-jordan.mock@example.com',
        'member',
        NOW() - INTERVAL '15 days',
        NOW(),
        'Jordan Lee',
        TRUE,
        NULL
    )
ON CONFLICT (id) DO UPDATE SET
    role = EXCLUDED.role,
    updated_at = NOW(),
    full_name = EXCLUDED.full_name,
    is_student = EXCLUDED.is_student,
    avatar_url = EXCLUDED.avatar_url;

-- Keep this seed rerunnable without deleting real audit logs.
DELETE FROM admin_audit_logs
WHERE request_id LIKE 'seed-admin-audit-%';

INSERT INTO admin_audit_logs (
    occurred_at,
    actor_user_id,
    action,
    target_user_id,
    outcome,
    request_id,
    description
)
VALUES
    (
        NOW() - INTERVAL '6 days',
        'a1000000-0000-4000-8000-000000000001',
        'users.viewed',
        NULL,
        'success',
        'seed-admin-audit-001',
        'Viewed the admin user directory'
    ),
    (
        NOW() - INTERVAL '5 days 12 hours',
        'a2000000-0000-4000-8000-000000000002',
        'users.exported',
        NULL,
        'success',
        'seed-admin-audit-002',
        'Exported the filtered user list as CSV'
    ),
    (
        NOW() - INTERVAL '4 days',
        'a1000000-0000-4000-8000-000000000001',
        'user.role_changed',
        'b1000000-0000-4000-8000-000000000001',
        'success',
        'seed-admin-audit-003',
        'Changed the user role from member to admin'
    ),
    (
        NOW() - INTERVAL '3 days 6 hours',
        'a2000000-0000-4000-8000-000000000002',
        'user.role_changed',
        'b2000000-0000-4000-8000-000000000002',
        'denied',
        'seed-admin-audit-004',
        'Role change was denied by authorization rules'
    ),
    (
        NOW() - INTERVAL '2 days',
        'a1000000-0000-4000-8000-000000000001',
        'user.group_assigned',
        'b1000000-0000-4000-8000-000000000001',
        'failed',
        'seed-admin-audit-005',
        'Could not assign the requested group'
    ),
    (
        NOW() - INTERVAL '2 hours',
        'a2000000-0000-4000-8000-000000000002',
        'users.viewed',
        NULL,
        'success',
        'seed-admin-audit-006',
        NULL
    );

COMMIT;
