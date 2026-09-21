-- ============================================================
-- MIGRATION: 004_create_request_status_logs_table.sql
-- Description: Audit trail — records every status change on a request
-- (BONUS: Audit Trail feature)
-- ============================================================

CREATE TABLE IF NOT EXISTS request_status_logs (
    id          UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id  UUID           NOT NULL REFERENCES maintenance_requests(id) ON DELETE CASCADE,
    changed_by  UUID           NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    from_status request_status,             -- NULL means the request was just created
    to_status   request_status NOT NULL,
    note        TEXT,                       -- Optional note for the status change
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- Index for fetching audit trail of a specific request
CREATE INDEX IF NOT EXISTS idx_status_logs_request_id
    ON request_status_logs(request_id, created_at DESC);

COMMENT ON TABLE  request_status_logs              IS 'Immutable audit log of all status changes for maintenance requests';
COMMENT ON COLUMN request_status_logs.from_status  IS 'Previous status; NULL when this log entry is the creation event';
COMMENT ON COLUMN request_status_logs.to_status    IS 'New status after the change';

