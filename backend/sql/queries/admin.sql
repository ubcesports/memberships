-- name: GetUsersAdmin :many
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
WHERE (
    sqlc.narg('full_name')::text IS NULL
    OR u.full_name ILIKE '%' || sqlc.narg('full_name')::text || '%'
)
AND (
    sqlc.narg('student_id')::text IS NULL
    OR u.student_id ILIKE '%' || sqlc.narg('student_id')::text || '%'
)
AND (
    sqlc.narg('email')::text IS NULL
    OR u.email ILIKE '%' || sqlc.narg('email')::text || '%'
)
AND (
    sqlc.narg('role')::role_type IS NULL
    OR u.role = sqlc.narg('role')::role_type
)
AND (
    sqlc.narg('is_student')::boolean IS NULL
    OR u.is_student = sqlc.narg('is_student')::boolean
)
AND (
    sqlc.narg('group')::group_type IS NULL
    OR EXISTS (
        SELECT 1
        FROM user_groups filter_group
        WHERE filter_group.user_id = u.id
          AND filter_group."group" = sqlc.narg('group')::group_type
    )
)
ORDER BY u.created_at DESC
LIMIT sqlc.narg('limit')
OFFSET sqlc.narg('offset');

-- name: CountUsersAdmin :one
SELECT COUNT(*)
FROM users u
WHERE (
    sqlc.narg('full_name')::text IS NULL
    OR u.full_name ILIKE '%' || sqlc.narg('full_name')::text || '%'
)
AND (
    sqlc.narg('student_id')::text IS NULL
    OR u.student_id ILIKE '%' || sqlc.narg('student_id')::text || '%'
)
AND (
    sqlc.narg('email')::text IS NULL
    OR u.email ILIKE '%' || sqlc.narg('email')::text || '%'
)
AND (
    sqlc.narg('role')::role_type IS NULL
    OR u.role = sqlc.narg('role')::role_type
)
AND (
    sqlc.narg('is_student')::boolean IS NULL
    OR u.is_student = sqlc.narg('is_student')::boolean
)
AND (
    sqlc.narg('group')::group_type IS NULL
    OR EXISTS (
        SELECT 1
        FROM user_groups filter_group
        WHERE filter_group.user_id = u.id
          AND filter_group."group" = sqlc.narg('group')::group_type
    )
);
