BEGIN;

INSERT INTO membership_tiers (
    title,
    slug,
    description,
    stripe_product_id,
    is_active
)
VALUES
    ('Day Pass', 'day', NULL, 'prod_UlEQ2pcezb8VkF', TRUE),
    ('Regular Pass', 'regular', NULL, 'prod_Ukn67a5I07Y1W4', TRUE),
    ('Premium Pass', 'premium', NULL, 'prod_Ukn7kJLDSdtB0t', TRUE)
ON CONFLICT (slug) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    stripe_product_id = EXCLUDED.stripe_product_id,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

DELETE FROM membership_tier_prices AS membership_price
USING membership_tiers AS membership_tier
WHERE membership_price.tier_id = membership_tier.id
  AND membership_tier.slug IN ('day', 'regular', 'premium')
  AND NOT EXISTS (
      SELECT 1
      FROM (
          VALUES
              ('day', 'member'),
              ('day', 'student'),
              ('regular', 'member'),
              ('regular', 'student'),
              ('regular', 'competitive_team'),
              ('regular', 'executive'),
              ('premium', 'member'),
              ('premium', 'student'),
              ('premium', 'competitive_team'),
              ('premium', 'executive')
      ) AS expected(tier_slug, price_group)
      WHERE expected.tier_slug = membership_tier.slug
        AND expected.price_group::group_type = membership_price."group"
  );

INSERT INTO membership_tier_prices (tier_id, "group", stripe_price_id)
SELECT
    membership_tier.id,
    price_mapping.price_group::group_type,
    price_mapping.stripe_price_id
FROM (
    VALUES
        ('day', 'member', 'price_1TlhxaDEhs9s474KjGKdv4Zl'),
        ('day', 'student', 'price_1TlhxGDEhs9s474KiqWxmkwk'),
        ('regular', 'member', 'price_1TlHWpDEhs9s474K9jOlZyQ5'),
        ('regular', 'student', 'price_1TlHW2DEhs9s474KMEjggTW8'),
        ('regular', 'competitive_team', 'price_1TlhtWDEhs9s474K5vSy8xZu'),
        ('regular', 'executive', 'price_1TlhthDEhs9s474KdEM5SlKc'),
        ('premium', 'member', 'price_1TlHXjDEhs9s474K4A6b4FrD'),
        ('premium', 'student', 'price_1TlHX8DEhs9s474KHf4xOEuz'),
        ('premium', 'competitive_team', 'price_1TlhwBDEhs9s474Kg8YlQw4k'),
        ('premium', 'executive', 'price_1TlhwLDEhs9s474Ki9HxkP8E')
) AS price_mapping(tier_slug, price_group, stripe_price_id)
JOIN membership_tiers AS membership_tier
  ON membership_tier.slug = price_mapping.tier_slug
ON CONFLICT (tier_id, "group") DO UPDATE SET
    stripe_price_id = EXCLUDED.stripe_price_id,
    updated_at = NOW();

COMMIT;
