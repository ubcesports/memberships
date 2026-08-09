-- name: GetUsersAdmin :many
WITH args AS (
    SELECT
        sqlc.narg('full_name')::text AS full_name,
        sqlc.narg('student_id')::text AS student_id,
        sqlc.narg('email')::text AS email,
        sqlc.narg('role')::role_type AS role,
        sqlc.narg('is_student')::boolean AS is_student,
        sqlc.narg('group')::group_type AS "group",
        sqlc.narg('limit')::integer AS "limit",
        sqlc.narg('offset')::integer AS "offset"
)
SELECT
    u.id,
    u.email,
    u.student_id,
    u.role,
    u.created_at,
    u.updated_at,
    u.full_name,
    u.email_verified_at,
    u.is_student,
    u.onboarding_completed_at,
    u.avatar_url,
    COALESCE(g.groups, '{}'::text[])::text[] AS groups
FROM users u
CROSS JOIN args a
LEFT JOIN LATERAL (
    SELECT array_agg(
        ug."group"::text
        ORDER BY ug.assigned_at, ug."group"
    ) AS groups
    FROM user_groups ug
    WHERE ug.user_id = u.id
) g ON true
WHERE (
    a.full_name IS NULL
    OR u.full_name ILIKE '%' || a.full_name || '%'
)
AND (
    a.student_id IS NULL
    OR u.student_id ILIKE '%' || a.student_id || '%'
)
AND (
    a.email IS NULL
    OR u.email ILIKE '%' || a.email || '%'
)
AND (
    a.role IS NULL
    OR u.role = a.role
)
AND (
    a.is_student IS NULL
    OR u.is_student = a.is_student
)
AND (
    a."group" IS NULL
    OR EXISTS (
        SELECT 1
        FROM user_groups filter_group
        WHERE filter_group.user_id = u.id
          AND filter_group."group" = a."group"
    )
)
ORDER BY u.created_at DESC
LIMIT (SELECT "limit" FROM args)
OFFSET (SELECT "offset" FROM args);

-- name: CountUsersAdmin :one
WITH args AS (
    SELECT
        sqlc.narg('full_name')::text AS full_name,
        sqlc.narg('student_id')::text AS student_id,
        sqlc.narg('email')::text AS email,
        sqlc.narg('role')::role_type AS role,
        sqlc.narg('is_student')::boolean AS is_student,
        sqlc.narg('group')::group_type AS "group"
)
SELECT COUNT(*)
FROM users u
CROSS JOIN args a
WHERE (
    a.full_name IS NULL
    OR u.full_name ILIKE '%' || a.full_name || '%'
)
AND (
    a.student_id IS NULL
    OR u.student_id ILIKE '%' || a.student_id || '%'
)
AND (
    a.email IS NULL
    OR u.email ILIKE '%' || a.email || '%'
)
AND (
    a.role IS NULL
    OR u.role = a.role
)
AND (
    a.is_student IS NULL
    OR u.is_student = a.is_student
)
AND (
    a."group" IS NULL
    OR EXISTS (
        SELECT 1
        FROM user_groups filter_group
        WHERE filter_group.user_id = u.id
          AND filter_group."group" = a."group"
    )
);

-- name: GetAdminUserByID :one
SELECT
    u.id,
    u.email,
    u.student_id,
    u.role,
    u.created_at,
    u.updated_at,
    u.full_name,
    u.email_verified_at,
    u.is_student,
    u.onboarding_completed_at,
    u.avatar_url,
    COALESCE(g.groups, '{}'::text[])::text[] AS groups
FROM users u
LEFT JOIN LATERAL (
    SELECT array_agg(
        ug."group"::text
        ORDER BY ug.assigned_at, ug."group"
    ) AS groups
    FROM user_groups ug
    WHERE ug.user_id = u.id
) g ON true
WHERE u.id = $1;

-- name: UpdateUserStudentInfo :exec
UPDATE users
SET
    is_student = $2,
    student_id = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: StudentIDExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE student_id = $1
);

-- name: UpdateUserRole :exec
UPDATE users
SET
    role = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: AddUserGroup :exec
INSERT INTO user_groups (user_id, "group")
VALUES ($1, $2)
ON CONFLICT (user_id, "group") DO NOTHING;

-- name: RemoveUserGroup :exec
DELETE FROM user_groups
WHERE user_id = $1 AND "group" = $2;

-- name: HasActiveMembershipForUser :one
SELECT EXISTS (
    SELECT 1
    FROM memberships
    WHERE user_id = $1
      AND cancelled_at IS NULL
);

-- name: GetMostRecentCancelledMembership :one
SELECT id, expires_at, cancelled_at
FROM memberships
WHERE user_id = $1
  AND cancelled_at IS NOT NULL
ORDER BY cancelled_at DESC
LIMIT 1;

-- name: ReinstateMembership :execrows
UPDATE memberships
SET
    cancelled_at = NULL,
    updated_at = NOW()
WHERE id = $1
  AND cancelled_at IS NOT NULL
  AND expires_at > NOW();

-- name: CreateAdminAuditLog :exec
INSERT INTO admin_audit_logs (
    actor_user_id,
    action,
    target_user_id,
    outcome,
    request_id,
    description
) VALUES (
    sqlc.arg('actor_user_id')::uuid,
    sqlc.arg('action')::text,
    sqlc.narg('target_user_id')::uuid,
    sqlc.arg('outcome')::admin_audit_outcome_type,
    sqlc.arg('request_id')::text,
    sqlc.narg('description')::text
);

-- name: GetAdminAuditLogs :many
WITH args AS (
    SELECT
        sqlc.narg('actor_name')::text AS actor_name,
        sqlc.arg('limit')::integer AS "limit",
        sqlc.arg('offset')::integer AS "offset"
)
SELECT
    aal.id,
    aal.occurred_at,
    aal.action,
    aal.outcome,
    aal.request_id,
    aal.description,
    actor.id AS actor_id,
    actor.full_name AS actor_name,
    actor.avatar_url AS actor_avatar_url,
    target.id AS target_id,
    target.full_name AS target_name,
    target.avatar_url AS target_avatar_url
FROM admin_audit_logs aal
JOIN users actor
    ON actor.id = aal.actor_user_id
LEFT JOIN users target
    ON target.id = aal.target_user_id
CROSS JOIN args a
WHERE (
    a.actor_name IS NULL
    OR actor.full_name ILIKE '%' || a.actor_name || '%'
)
ORDER BY
    aal.occurred_at DESC,
    aal.id DESC
LIMIT (SELECT "limit" FROM args)
OFFSET (SELECT "offset" FROM args);
