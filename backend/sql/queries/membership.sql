-- name: ListEligibleTierPriceMappings :many
SELECT
    mt.id AS tier_id,
    mt.slug,
    mt.title,
    mt.description,
    mt.stripe_product_id,
    mtp."group",
    mtp.is_student,
    mtp.stripe_price_id
FROM users u
JOIN membership_tier_prices mtp
    ON mtp.is_student = u.is_student
   AND (
       mtp."group" = 'member'
       OR EXISTS (
           SELECT 1
           FROM user_groups ug
           WHERE ug.user_id = u.id
             AND ug."group" = mtp."group"
       )
   )
JOIN membership_tiers mt ON mt.id = mtp.tier_id
WHERE u.id = $1
  AND mt.is_active = TRUE
  AND mt.stripe_product_id IS NOT NULL
ORDER BY mt.title, mtp."group";

-- name: GetActiveMembershipByUserID :one
SELECT
    m.id,
    m.user_id,
    m.tier_id,
    mt.slug AS tier_slug,
    mt.title AS tier_title,
    mt.description AS tier_description,
    m.group_at_purchase,
    m.is_student_at_purchase,
    m.started_at,
    m.expires_at,
    m.cancelled_at,
    m.created_at,
    m.updated_at
FROM memberships m
JOIN membership_tiers mt ON mt.id = m.tier_id
WHERE m.user_id = $1
  AND m.cancelled_at IS NULL
  AND m.expires_at > NOW()
ORDER BY m.started_at DESC
LIMIT 1;

-- name: LockUserForCheckout :one
SELECT id
FROM users
WHERE id = $1
FOR UPDATE;

-- name: GetUserEmail :one
SELECT email
FROM users
WHERE id = $1;

-- name: GetPendingTransactionByUserID :one
SELECT *
FROM transactions
WHERE user_id = $1
  AND status = 'pending'
LIMIT 1;

-- name: GetTransactionBySessionIDForUpdate :one
SELECT *
FROM transactions
WHERE stripe_checkout_session_id = $1
FOR UPDATE;

-- name: GetTransactionByPaymentIntentForUpdate :one
SELECT *
FROM transactions
WHERE stripe_payment_intent_id = $1
FOR UPDATE;

-- name: GetTransactionByIDForUpdate :one
SELECT *
FROM transactions
WHERE id = $1
FOR UPDATE;

-- name: CreatePendingTransaction :one
INSERT INTO transactions (
    id,
    user_id,
    tier_id,
    group_at_purchase,
    is_student_at_purchase,
    stripe_price_id,
    amount_minor,
    currency,
    status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
RETURNING *;

-- name: AttachCheckoutSession :exec
UPDATE transactions
SET stripe_checkout_session_id = $2, updated_at = NOW()
WHERE id = $1
  AND status = 'pending';

-- name: MarkPendingTransactionFailed :execrows
UPDATE transactions
SET status = 'failed', updated_at = NOW()
WHERE id = $1
  AND status = 'pending';

-- name: MarkPendingTransactionExpired :execrows
UPDATE transactions
SET status = 'expired', updated_at = NOW()
WHERE id = $1
  AND status = 'pending';

-- name: CreateMembership :one
INSERT INTO memberships (
    id,
    user_id,
    tier_id,
    group_at_purchase,
    is_student_at_purchase,
    started_at,
    expires_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: CompleteTransaction :exec
UPDATE transactions
SET
    membership_id = $2,
    stripe_payment_intent_id = $3,
    stripe_charge_id = $4,
    amount_minor = $5,
    currency = $6,
    status = 'completed',
    updated_at = NOW()
WHERE id = $1
  AND status IN ('pending', 'failed');

-- name: MarkTransactionRefunded :exec
UPDATE transactions
SET
    stripe_payment_intent_id = COALESCE(NULLIF($2::varchar, ''), stripe_payment_intent_id),
    stripe_charge_id = COALESCE(NULLIF($3::varchar, ''), stripe_charge_id),
    status = 'refunded',
    updated_at = NOW()
WHERE id = $1
  AND status <> 'refunded';

-- name: CancelMembership :exec
UPDATE memberships
SET cancelled_at = COALESCE(cancelled_at, $2), updated_at = NOW()
WHERE id = $1;
