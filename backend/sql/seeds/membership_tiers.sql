BEGIN;

WITH tier_seed(title, description, slug, group_name, stripe_product_id, is_active) AS (
    VALUES
        (
            'Day Pass',
            'Single-day UBC Esports Association membership access.',
            'day',
            'member',
            'prod_UlEQ2pcezb8VkF',
            TRUE
        ),
        (
            'Regular Pass',
            'Standard UBC Esports Association membership pass.',
            'regular',
            'member',
            'prod_UnGYL5Cxq8Sr7r',
            TRUE
        ),
        (
            'Premium Pass',
            'Premium UBC Esports Association membership pass.',
            'premium',
            'member',
            'prod_UnGYUlfwqi6af8',
            TRUE
        ),
        (
            'Competitive Player Pass',
            'Membership pass for competitive team players.',
            'competitive_team',
            'competitive_team',
            'prod_UnGZnjFGP74y44',
            TRUE
        ),
        (
            'Executive Pass',
            'Membership pass for executive, director, and board members.',
            'executive',
            'executive',
            'prod_UnGZIuiCafLoqL',
            TRUE
        )
),
upserted_tiers AS (
    INSERT INTO membership_tiers (
        title,
        description,
        slug,
        "group",
        stripe_product_id,
        is_active,
        updated_at
    )
    SELECT
        title,
        description,
        slug,
        group_name::group_type,
        stripe_product_id,
        is_active,
        NOW()
    FROM tier_seed
    ON CONFLICT (stripe_product_id) DO UPDATE SET
        title = EXCLUDED.title,
        description = EXCLUDED.description,
        slug = EXCLUDED.slug,
        "group" = EXCLUDED."group",
        is_active = EXCLUDED.is_active,
        updated_at = NOW()
    RETURNING id, slug
),
price_seed(slug, stripe_price_id, is_student_required) AS (
    VALUES
        ('day', 'price_1TlhxGDEhs9s474KiqWxmkwk', TRUE),
        ('day', 'price_1TlhxaDEhs9s474KjGKdv4Zl', FALSE),
        ('regular', 'price_1Tng1SDEhs9s474KmB0chFf4', TRUE),
        ('regular', 'price_1Tng2XDEhs9s474KeFwj2xfZ', FALSE),
        ('premium', 'price_1Tng1hDEhs9s474KJUoSaOy8', TRUE),
        ('premium', 'price_1Tng30DEhs9s474KWdiXE6iK', FALSE),
        ('competitive_team', 'price_1Tng27DEhs9s474KStx6hYtV', NULL),
        ('executive', 'price_1Tng1sDEhs9s474KxHto1E80', TRUE),
        ('executive', 'price_1Tng3TDEhs9s474K0MpRQlaG', FALSE)
)
INSERT INTO membership_tier_prices (
    tier_id,
    stripe_price_id,
    is_student_required,
    updated_at
)
SELECT
    t.id,
    p.stripe_price_id,
    p.is_student_required,
    NOW()
FROM price_seed p
JOIN upserted_tiers t
    ON t.slug = p.slug
ON CONFLICT (stripe_price_id) DO UPDATE SET
    tier_id = EXCLUDED.tier_id,
    is_student_required = EXCLUDED.is_student_required,
    updated_at = NOW();

COMMIT;
