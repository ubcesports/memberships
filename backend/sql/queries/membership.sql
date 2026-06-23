-- name: ListPublicTierPriceMappings :many
SELECT
    mt.id AS tier_id,
    mt.slug,
    mt.title,
    mt.description,
    mt.stripe_product_id,
    mt.is_public,
    mt.required_group,
    mtp."group",
    mtp.stripe_price_id
FROM membership_tiers mt
JOIN membership_tier_prices mtp ON mtp.tier_id = mt.id
WHERE mt.is_active = TRUE
  AND mt.is_public = TRUE
  AND mt.stripe_product_id IS NOT NULL
  AND mtp."group" IN ('member', 'student')
ORDER BY mt.title, mtp."group";

-- name: ListEligibleTierPriceMappings :many
SELECT
    mt.id AS tier_id,
    mt.slug,
    mt.title,
    mt.description,
    mt.stripe_product_id,
    mt.is_public,
    mt.required_group,
    mtp."group",
    mtp.stripe_price_id
FROM membership_tiers mt
JOIN membership_tier_prices mtp ON mtp.tier_id = mt.id
WHERE mt.is_active = TRUE
  AND mt.stripe_product_id IS NOT NULL
  AND (
      (
          mt.required_group IS NULL
          AND (
              (
                  mtp."group" = 'student'
                  AND EXISTS (
                      SELECT 1 FROM user_groups ug
                      WHERE ug.user_id = $1 AND ug."group" = 'student'
                  )
              )
              OR (
                  mtp."group" = 'member'
                  AND NOT EXISTS (
                      SELECT 1 FROM user_groups ug
                      WHERE ug.user_id = $1 AND ug."group" = 'student'
                  )
              )
          )
      )
      OR (
          mt.required_group IS NOT NULL
          AND mtp."group" = mt.required_group
          AND EXISTS (
              SELECT 1 FROM user_groups ug
              WHERE ug.user_id = $1 AND ug."group" = mt.required_group
          )
      )
  )
ORDER BY mt.title;

-- name: GetActiveMembershipByUserID :one
SELECT
    m.id,
    m.user_id,
    m.tier_id,
    mt.slug AS tier_slug,
    mt.title AS tier_title,
    mt.description AS tier_description,
    m.group_at_purchase,
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

-- name: GetCompletedPaidAmountForMembership :one
SELECT COALESCE(SUM(amount_minor), 0)::bigint
FROM transactions
WHERE membership_id = $1
  AND status = 'completed';

-- name: LockUserForCheckout :one
SELECT id FROM users WHERE id = $1 FOR UPDATE;

-- name: GetUserEmail :one
SELECT email FROM users WHERE id = $1;

-- name: GetPendingTransactionByUserID :one
SELECT * FROM transactions
WHERE user_id = $1 AND status = 'pending'
LIMIT 1;

-- name: GetTransactionBySessionIDForUpdate :one
SELECT * FROM transactions
WHERE stripe_checkout_session_id = $1
FOR UPDATE;

-- name: CreatePendingTransaction :one
INSERT INTO transactions (
    id,
    user_id,
    membership_id,
    tier_id,
    group_at_purchase,
    stripe_price_id,
    amount_minor,
    credit_amount_minor,
    currency,
    kind,
    status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'pending')
RETURNING *;

-- name: AttachCheckoutSession :exec
UPDATE transactions
SET stripe_checkout_session_id = $2, updated_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: MarkPendingTransactionFailed :execrows
UPDATE transactions
SET status = 'failed', updated_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: MarkPendingTransactionExpired :execrows
UPDATE transactions
SET status = 'expired', updated_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: CreateMembership :one
INSERT INTO memberships (
    id,
    user_id,
    tier_id,
    group_at_purchase,
    started_at,
    expires_at
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CancelMembershipForUpgrade :one
UPDATE memberships
SET
    cancelled_at = $2,
    updated_at = NOW()
WHERE memberships.id = $1
  AND cancelled_at IS NULL
  AND expires_at > NOW()
  AND EXISTS (
      SELECT 1
      FROM membership_tiers current_tier
      WHERE current_tier.id = memberships.tier_id
        AND current_tier.slug = 'regular'
  )
RETURNING expires_at;

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
