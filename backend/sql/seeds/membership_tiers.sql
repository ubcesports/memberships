BEGIN;

INSERT INTO membership_tiers (title, slug, description, is_active, is_public, required_group)
VALUES
    ('Regular', 'regular', 'Regular UBCEA membership pass', TRUE, TRUE, NULL),
    ('Premium', 'premium', 'Premium UBCEA membership pass', TRUE, TRUE, NULL),
    ('Competitive Team', 'competitive-team', 'Competitive Team UBCEA membership pass', TRUE, FALSE, 'competitive_team'),
    ('Executive', 'executive', 'Executive UBCEA membership pass', TRUE, FALSE, 'executive')
ON CONFLICT (slug) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    is_active = EXCLUDED.is_active,
    is_public = EXCLUDED.is_public,
    required_group = EXCLUDED.required_group,
    updated_at = NOW();

COMMIT;
