-- name: CompleteUserOnboarding :exec
UPDATE users
SET is_student = $2,
    student_id = $3,
    onboarding_completed_at = NOW()
WHERE id = $1;
