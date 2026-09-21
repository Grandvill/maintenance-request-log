-- ============================================================
-- MIGRATION: 003_create_maintenance_requests_table.sql
-- Description: Core maintenance request table with FK to users
-- ============================================================

CREATE TABLE IF NOT EXISTS maintenance_requests (
    id                  UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id            VARCHAR(50)    NOT NULL,
    problem_description TEXT           NOT NULL,
    priority            request_priority NOT NULL DEFAULT 'Medium',
    status              request_status   NOT NULL DEFAULT 'Submitted',

    -- FK: who created this request
    created_by          UUID           NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    -- FK: who reviewed (approved/rejected) this request
    reviewed_by         UUID           REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at         TIMESTAMPTZ,
    review_note         TEXT,           -- Optional note/reason from reviewer

    created_at          TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- Index for operator filtering: "show only my requests"
CREATE INDEX IF NOT EXISTS idx_requests_created_by
    ON maintenance_requests(created_by);

-- Index for status filtering (most common query)
CREATE INDEX IF NOT EXISTS idx_requests_status
    ON maintenance_requests(status);

-- Index for priority filtering
CREATE INDEX IF NOT EXISTS idx_requests_priority
    ON maintenance_requests(priority);

-- Composite index for filtered list (status + created_at for sorting)
CREATE INDEX IF NOT EXISTS idx_requests_status_created
    ON maintenance_requests(status, created_at DESC);

-- Comments
COMMENT ON TABLE  maintenance_requests                  IS 'Maintenance requests submitted by operators';
COMMENT ON COLUMN maintenance_requests.asset_id         IS 'Identifier of the asset/machine that needs maintenance';
COMMENT ON COLUMN maintenance_requests.review_note      IS 'Optional reason/comment when supervisor approves or rejects';
COMMENT ON COLUMN maintenance_requests.reviewed_at      IS 'Timestamp when the request was approved or rejected';

