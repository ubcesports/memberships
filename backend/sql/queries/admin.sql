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
