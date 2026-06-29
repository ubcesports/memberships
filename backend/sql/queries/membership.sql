-- name: GetPublicTiersAndPrices :many
SELECT
    mt.id,
    mt.title,
    mt.description,
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
