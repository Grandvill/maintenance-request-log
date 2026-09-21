-- ============================================================
-- MIGRATION: 005_create_updated_at_trigger.sql
-- Description: Auto-update updated_at column on any row UPDATE
-- ============================================================

-- Function used by all triggers
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger on users table
DROP TRIGGER IF EXISTS set_users_updated_at ON users;
CREATE TRIGGER set_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Trigger on maintenance_requests table
DROP TRIGGER IF EXISTS set_requests_updated_at ON maintenance_requests;
CREATE TRIGGER set_requests_updated_at
    BEFORE UPDATE ON maintenance_requests
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

