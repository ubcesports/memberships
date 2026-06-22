-- name: GetProfileByUserID :one
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
    SELECT array_agg(ug."group"::text ORDER BY ug.assigned_at ASC, ug."group" ASC) AS groups
    FROM user_groups ug
    WHERE ug.user_id = u.id
) g ON true
WHERE u.id = $1
;
