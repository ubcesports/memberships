-- Seed membership tiers and prices.
-- Run after migrations, for example:
-- psql "$DATABASE_URL" -f backend/sql/seeds/membership_tiers_and_prices.sql

INSERT INTO membership_tiers (code, title, description)
VALUES
    ('regular', 'Regular Pass', 'One 1-hour Lounge session per day, Chunithm at $1 per 3 songs, and UBC Smash Weekly discount.'),
    ('premium', 'Premium Pass', 'Unlimited 2-hour Lounge sessions per day when space is available, Chunithm access, and UBC Smash Weekly discount.'),
    ('cab', 'Cab Pass', 'Chunithm cab access only.'),
    ('day', 'Day Pass', 'One-day Lounge access.')
ON CONFLICT (code) DO UPDATE
SET title = EXCLUDED.title,
    description = EXCLUDED.description,
    updated_at = NOW();

WITH tier_prices(code, "group", student_status, price) AS (
    VALUES
        ('regular'::tier_code_type, 'member'::group_type, 'student'::student_status_type, 15.00),
        ('regular'::tier_code_type, 'member'::group_type, 'non_student'::student_status_type, 20.00),
        ('premium'::tier_code_type, 'member'::group_type, 'student'::student_status_type, 25.00),
        ('premium'::tier_code_type, 'member'::group_type, 'non_student'::student_status_type, 30.00),
        ('cab'::tier_code_type, 'member'::group_type, 'student'::student_status_type, 10.00),
        ('cab'::tier_code_type, 'member'::group_type, 'non_student'::student_status_type, 15.00),
        ('day'::tier_code_type, 'member'::group_type, 'student'::student_status_type, 5.00),
        ('day'::tier_code_type, 'member'::group_type, 'non_student'::student_status_type, 10.00),

        ('regular'::tier_code_type, 'competitive_team'::group_type, 'student'::student_status_type, 0.00),
        ('regular'::tier_code_type, 'competitive_team'::group_type, 'non_student'::student_status_type, 0.00),
        ('premium'::tier_code_type, 'competitive_team'::group_type, 'student'::student_status_type, 0.00),
        ('premium'::tier_code_type, 'competitive_team'::group_type, 'non_student'::student_status_type, 0.00),
        ('cab'::tier_code_type, 'competitive_team'::group_type, 'student'::student_status_type, 10.00),
        ('cab'::tier_code_type, 'competitive_team'::group_type, 'non_student'::student_status_type, 15.00),
        ('day'::tier_code_type, 'competitive_team'::group_type, 'student'::student_status_type, 5.00),
        ('day'::tier_code_type, 'competitive_team'::group_type, 'non_student'::student_status_type, 10.00),

        ('regular'::tier_code_type, 'executive'::group_type, 'student'::student_status_type, 15.00),
        ('regular'::tier_code_type, 'executive'::group_type, 'non_student'::student_status_type, 20.00),
        ('premium'::tier_code_type, 'executive'::group_type, 'student'::student_status_type, 25.00),
        ('premium'::tier_code_type, 'executive'::group_type, 'non_student'::student_status_type, 30.00),
        ('cab'::tier_code_type, 'executive'::group_type, 'student'::student_status_type, 10.00),
        ('cab'::tier_code_type, 'executive'::group_type, 'non_student'::student_status_type, 15.00),
        ('day'::tier_code_type, 'executive'::group_type, 'student'::student_status_type, 5.00),
        ('day'::tier_code_type, 'executive'::group_type, 'non_student'::student_status_type, 10.00),

        ('regular'::tier_code_type, 'director'::group_type, 'student'::student_status_type, 15.00),
        ('regular'::tier_code_type, 'director'::group_type, 'non_student'::student_status_type, 20.00),
        ('premium'::tier_code_type, 'director'::group_type, 'student'::student_status_type, 25.00),
        ('premium'::tier_code_type, 'director'::group_type, 'non_student'::student_status_type, 30.00),
        ('cab'::tier_code_type, 'director'::group_type, 'student'::student_status_type, 10.00),
        ('cab'::tier_code_type, 'director'::group_type, 'non_student'::student_status_type, 15.00),
        ('day'::tier_code_type, 'director'::group_type, 'student'::student_status_type, 5.00),
        ('day'::tier_code_type, 'director'::group_type, 'non_student'::student_status_type, 10.00),

        ('regular'::tier_code_type, 'board'::group_type, 'student'::student_status_type, 15.00),
        ('regular'::tier_code_type, 'board'::group_type, 'non_student'::student_status_type, 20.00),
        ('premium'::tier_code_type, 'board'::group_type, 'student'::student_status_type, 25.00),
        ('premium'::tier_code_type, 'board'::group_type, 'non_student'::student_status_type, 30.00),
        ('cab'::tier_code_type, 'board'::group_type, 'student'::student_status_type, 10.00),
        ('cab'::tier_code_type, 'board'::group_type, 'non_student'::student_status_type, 15.00),
        ('day'::tier_code_type, 'board'::group_type, 'student'::student_status_type, 5.00),
        ('day'::tier_code_type, 'board'::group_type, 'non_student'::student_status_type, 10.00)
)
INSERT INTO membership_tier_prices (tier_id, "group", student_status, price)
SELECT mt.id, tier_prices."group", tier_prices.student_status, tier_prices.price
FROM tier_prices
JOIN membership_tiers mt ON mt.code = tier_prices.code
ON CONFLICT (tier_id, "group", student_status) DO UPDATE
SET price = EXCLUDED.price,
    updated_at = NOW();
