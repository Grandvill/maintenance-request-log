-- ============================================================
-- MIGRATION: 001_create_enum_types.sql
-- Description: Create all ENUM types used across tables
-- ============================================================

-- User roles
DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('operator', 'supervisor', 'admin');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Maintenance request status
DO $$ BEGIN
    CREATE TYPE request_status AS ENUM ('Submitted', 'Approved', 'Rejected');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Maintenance request priority
DO $$ BEGIN
    CREATE TYPE request_priority AS ENUM ('Low', 'Medium', 'High', 'Urgent');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

