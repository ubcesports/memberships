BEGIN;

INSERT INTO membership_tiers (title, slug, description, is_active)
VALUES
    ('Regular', 'regular', 'Regular UBCEA membership pass', TRUE),
    ('Premium', 'premium', 'Premium UBCEA membership pass', TRUE)
ON CONFLICT (slug) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

COMMIT;
