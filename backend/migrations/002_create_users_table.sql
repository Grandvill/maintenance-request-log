-- ============================================================
-- MIGRATION: 002_create_users_table.sql
-- Description: Create users table for authentication & RBAC
-- ============================================================

CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username    VARCHAR(50)  NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,   -- bcrypt hashed
    full_name   VARCHAR(100) NOT NULL,
    role        user_role    NOT NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Index for fast login lookup
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Index for filtering by role and active status
CREATE INDEX IF NOT EXISTS idx_users_role_active ON users(role, is_active);

-- Comments
COMMENT ON TABLE  users              IS 'Application users with role-based access control';
COMMENT ON COLUMN users.password     IS 'bcrypt hashed password (cost 12)';
COMMENT ON COLUMN users.role         IS 'operator | supervisor | admin';
COMMENT ON COLUMN users.is_active    IS 'Soft delete: false means deactivated, not deleted';

