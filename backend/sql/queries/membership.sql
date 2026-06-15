-- name: GetUserForMembershipPricing :one
SELECT id, is_student
FROM users
WHERE id = $1;

-- name: ListUserGroups :many
SELECT "group"
FROM user_groups
WHERE user_id = $1;

-- name: ListActiveMembershipTiersWithPrices :many
SELECT
    mt.id AS tier_id,
    mt.code,
    mt.title,
    mt.description,
    mtp.id AS price_id,
    mtp."group",
    mtp.student_status,
    mtp.price,
    mtp.stripe_price_id
FROM membership_tiers mt
JOIN membership_tier_prices mtp ON mtp.tier_id = mt.id
WHERE mt.is_active = TRUE
    AND mtp."group" = $1
    AND mtp.student_status = $2
ORDER BY
    CASE mt.code
        WHEN 'regular' THEN 1
        WHEN 'premium' THEN 2
        WHEN 'cab' THEN 3
        WHEN 'day' THEN 4
        ELSE 5
    END,
    mt.title;

-- name: GetTierPriceByCode :one
SELECT
    mt.id AS tier_id,
    mt.code,
    mt.title,
    mt.description,
    mtp.id AS price_id,
    mtp."group",
    mtp.student_status,
    mtp.price,
    mtp.stripe_price_id
FROM membership_tiers mt
JOIN membership_tier_prices mtp ON mtp.tier_id = mt.id
WHERE mt.is_active = TRUE
    AND mt.code = $1
    AND mtp."group" = $2
    AND mtp.student_status = $3;

-- name: GetCurrentMembershipByUserID :one
SELECT
    m.id AS membership_id,
    m.user_id,
    m.tier_id,
    mt.code,
    mt.title,
    mt.description,
    m.transaction_id,
    m.group_at_purchase,
    m.student_status_at_purchase,
    m.started_at,
    m.expires_at,
    m.cancelled_at,
    t.price_amount,
    t.status AS transaction_status,
    t.stripe_payment_intent_id
FROM memberships m
JOIN membership_tiers mt ON mt.id = m.tier_id
LEFT JOIN transactions t ON t.id = m.transaction_id
WHERE m.user_id = $1
    AND m.cancelled_at IS NULL;

-- name: CreateMembership :one
INSERT INTO memberships (
    user_id,
    tier_id,
    group_at_purchase,
    student_status_at_purchase,
    started_at,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id,
    stripe_payment_intent_id,
    price_amount,
    status
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: SetMembershipTransactionID :one
UPDATE memberships
SET transaction_id = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetTransactionByStripePaymentIntentID :one
SELECT *
FROM transactions
WHERE stripe_payment_intent_id = $1;

-- name: UpdateTransactionStatus :one
UPDATE transactions
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
