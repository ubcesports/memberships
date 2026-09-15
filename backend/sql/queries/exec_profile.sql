-- name: GetExecProfiles :many
SELECT 
    e.user_id,
    e.title,
    e.display_order,
    e.display_group,
    u.full_name,
    u.avatar_url,
    COALESCE(
        json_agg(
            json_build_object(
                'platform', s.platform,
                'url', s.url
            )
        ), '[]'
    )::json AS social_links
FROM exec_profile e
JOIN users u ON e.user_id = u.id
LEFT JOIN exec_social_link s ON e.user_id = s.user_id
GROUP BY e.user_id, e.title, e.display_order, e.display_group, u.full_name, u.avatar_url
ORDER BY e.display_order ASC;

-- name: AddExecSocialLink :exec
INSERT INTO exec_social_link (user_id, platform, url)
VALUES ($1, $2, $3);

-- name: UpdateExecSocialLink :exec
UPDATE exec_social_link
SET url = $3, updated_at = NOW()
WHERE user_id = $1 AND platform = $2;

-- name: DeleteExecSocialLink :exec
DELETE FROM exec_social_link 
WHERE user_id = $1 AND platform = $2;

-- name: UpdateExecProfileTitle :exec
UPDATE exec_profile
SET title = $2, updated_at = NOW()
WHERE user_id = $1;