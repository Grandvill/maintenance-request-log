-- ============================================================
-- SEED: 001_seed_users.sql
-- Description: Insert default system users for all 3 roles.
--
-- Credentials (for README / testing):
--   admin      / admin123
--   supervisor / supervisor123
--   operator1  / operator123
--   operator2  / operator123
--
-- Passwords are bcrypt hashed with cost=12.
-- Hash generated via: https://bcrypt-generator.com/ or Go's bcrypt package.
-- ============================================================

-- NOTE: This script uses INSERT ... ON CONFLICT DO NOTHING
-- so it is safe to re-run without duplicating data.

INSERT INTO users (id, username, password, full_name, role, is_active)
VALUES
    -- Admin user
    (
        '00000000-0000-0000-0000-000000000001',
        'admin',
        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/RKHm5yGji', -- admin123
        'System Administrator',
        'admin',
        TRUE
    ),
    -- Supervisor user
    (
        '00000000-0000-0000-0000-000000000002',
        'supervisor',
        '$2a$12$5R.7HQVPKzRdXbpxHFnv0uJV7s0b1dC3n6f8Ky2.wqA4H5x3mEJvO', -- supervisor123
        'Budi Santoso',
        'supervisor',
        TRUE
    ),
    -- Operator users
    (
        '00000000-0000-0000-0000-000000000003',
        'operator1',
        '$2a$12$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- operator123
        'Andi Prasetyo',
        'operator',
        TRUE
    ),
    (
        '00000000-0000-0000-0000-000000000004',
        'operator2',
        '$2a$12$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- operator123
        'Siti Rahayu',
        'operator',
        TRUE
    ),
    -- Inactive operator (to demonstrate deactivation feature)
    (
        '00000000-0000-0000-0000-000000000005',
        'operator_inactive',
        '$2a$12$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- operator123
        'Joko Widodo',
        'operator',
        FALSE
    )
ON CONFLICT (id) DO NOTHING;

