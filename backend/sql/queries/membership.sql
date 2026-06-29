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