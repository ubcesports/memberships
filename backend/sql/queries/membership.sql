-- name: GetPublicTiersAndPrices :many
SELECT
    mt.id,
    mt.title,
    mt.description,
    mt.benefits,
    mt.slug,
    mt.stripe_product_id,
    mtp.stripe_price_id,
    mtp.is_student_required
FROM membership_tiers mt
JOIN membership_tier_prices mtp
    ON mtp.tier_id = mt.id
WHERE mt.is_active = TRUE AND mt."group" = 'member';

-- name: GetCurrentMembershipWithTransaction :one
SELECT
    m.id,
    m.tier_id,
    m.started_at,
    m.expires_at,
    m.cancelled_at,
    t.id AS transaction_id,
    t.amount_paid_cents,
    t.status,
    t.group_at_purchase
FROM memberships m
JOIN transactions t
    ON t.membership_id = m.id
WHERE m.user_id = $1
    AND m.cancelled_at IS NULL
    AND m.started_at <= NOW()
    AND m.expires_at > NOW()
ORDER BY m.started_at DESC
LIMIT 1;

-- name: GetAllMembershipsWithTransactions :many
SELECT
    m.id,
    m.tier_id,
    m.started_at,
    m.expires_at,
    m.cancelled_at,
    t.id AS transaction_id,
    t.amount_paid_cents,
    t.status,
    t.group_at_purchase
FROM memberships m
JOIN transactions t
    ON t.membership_id = m.id
WHERE m.user_id = $1
ORDER BY m.started_at DESC;

-- name: GetEligibleTiersWithPrices :many
SELECT
    mt.id,
    mt.title,
    mt.description,
    mt.benefits,
    mt.stripe_product_id,
    mt.slug,
    mtp.stripe_price_id,
    mtp.is_student_required
FROM membership_tiers mt
JOIN membership_tier_prices mtp
    ON mtp.tier_id = mt.id
JOIN users u
    on u.id = $1
WHERE mt.is_active = TRUE
    AND (
        mt."group" = 'member'
        OR EXISTS (
            SELECT 1
            FROM user_groups ug
            WHERE ug.user_id = u.id AND ug."group" = mt."group"
        )
    )
    AND (
        mtp.is_student_required IS NULL
        OR mtp.is_student_required = u.is_student
    );

-- name: GetTierByTierId :one
SELECT
    mt.id,
    mt.title,
    mt.description,
    mt.benefits,
    mt.slug,
    mt.stripe_product_id,
    mtp.stripe_price_id,
    mtp.is_student_required
FROM membership_tiers mt
JOIN membership_tier_prices mtp
    ON mtp.tier_id = mt.id
WHERE mt.id = $1;

-- name: GetPendingTransactionForUpdate :one
SELECT
    id,
    stripe_checkout_session_id
FROM transactions
WHERE user_id = $1 AND status = 'pending'
FOR UPDATE;

-- name: ExpirePendingTransactionById :exec
UPDATE transactions
SET
    status = 'expired',
    updated_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: CreatePendingTransaction :one
INSERT INTO transactions (
    user_id,
    tier_id,
    group_at_purchase,
    student_at_purchase,
    purchase_type,
    status
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    'pending'
)
RETURNING id;

-- name: PutStripeCheckoutSessionId :exec
UPDATE transactions
SET
    stripe_checkout_session_id = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateTransactionStatusById :exec
UPDATE transactions
SET
    status = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdatePendingTransactionStatusByCheckoutId :exec
UPDATE transactions
SET
    status = $2,
    updated_at = NOW()
WHERE stripe_checkout_session_id = $1 AND status = 'pending';

-- name: GetTransactionByCheckoutSessionIdForUpdate :one
SELECT
    id,
    user_id,
    membership_id,
    tier_id,
    status,
    purchase_type,
    stripe_checkout_session_id
FROM transactions
WHERE stripe_checkout_session_id = $1
FOR UPDATE;

-- name: CreateMembership :one
INSERT INTO memberships (
    user_id,
    tier_id,
    started_at,
    expires_at
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id;

-- name: CompleteTransaction :exec
UPDATE transactions
SET
    membership_id = $2,
    stripe_payment_intent_id = $3,
    amount_paid_cents = $4,
    status = 'completed',
    updated_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: CancelActiveMembershipsByUserId :exec
UPDATE memberships
SET
    cancelled_at = $2,
    updated_at = $2
WHERE user_id = $1 AND cancelled_at IS NULL;