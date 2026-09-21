package requests

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// List queries requests with optional filters and join on users table
func (r *Repository) List(params FilterParams) ([]MaintenanceRequest, error) {
	query := `
		SELECT 
			r.id, r.asset_id, r.problem_description, r.priority, r.status,
			r.created_by, u_creator.full_name AS creator_name,
			r.reviewed_by, u_reviewer.full_name AS reviewer_name,
			r.reviewed_at, r.review_note,
			r.created_at, r.updated_at
		FROM maintenance_requests r
		JOIN users u_creator ON r.created_by = u_creator.id
		LEFT JOIN users u_reviewer ON r.reviewed_by = u_reviewer.id
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if params.CreatedBy != "" {
		query += fmt.Sprintf(" AND r.created_by = $%d", argIdx)
		args = append(args, params.CreatedBy)
		argIdx++
	}

	if params.Status != "" {
		query += fmt.Sprintf(" AND r.status = $%d", argIdx)
		args = append(args, params.Status)
		argIdx++
	}

	if params.Priority != "" {
		query += fmt.Sprintf(" AND r.priority = $%d", argIdx)
		args = append(args, params.Priority)
		argIdx++
	}

	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query += fmt.Sprintf(" AND (LOWER(r.asset_id) LIKE $%d OR LOWER(r.problem_description) LIKE $%d)", argIdx, argIdx)
		args = append(args, searchPattern)
		argIdx++
	}

	query += " ORDER BY r.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query requests: %w", err)
	}
	defer rows.Close()

	var result []MaintenanceRequest
	for rows.Next() {
		var req MaintenanceRequest
		err := rows.Scan(
			&req.ID, &req.AssetID, &req.ProblemDescription, &req.Priority, &req.Status,
			&req.CreatedBy, &req.CreatorName,
			&req.ReviewedBy, &req.ReviewerName,
			&req.ReviewedAt, &req.ReviewNote,
			&req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, req)
	}

	if result == nil {
		result = []MaintenanceRequest{}
	}

	return result, nil
}

// GetByID fetches a single request with creator and reviewer details
func (r *Repository) GetByID(id string) (*MaintenanceRequest, error) {
	query := `
		SELECT 
			r.id, r.asset_id, r.problem_description, r.priority, r.status,
			r.created_by, u_creator.full_name AS creator_name,
			r.reviewed_by, u_reviewer.full_name AS reviewer_name,
			r.reviewed_at, r.review_note,
			r.created_at, r.updated_at
		FROM maintenance_requests r
		JOIN users u_creator ON r.created_by = u_creator.id
		LEFT JOIN users u_reviewer ON r.reviewed_by = u_reviewer.id
		WHERE r.id = $1
	`
	var req MaintenanceRequest
	err := r.db.QueryRow(query, id).Scan(
		&req.ID, &req.AssetID, &req.ProblemDescription, &req.Priority, &req.Status,
		&req.CreatedBy, &req.CreatorName,
		&req.ReviewedBy, &req.ReviewerName,
		&req.ReviewedAt, &req.ReviewNote,
		&req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("maintenance request not found")
		}
		return nil, err
	}
	return &req, nil
}

// GetAuditLogs fetches status logs for a specific request
func (r *Repository) GetAuditLogs(requestID string) ([]StatusLog, error) {
	query := `
		SELECT 
			l.id, l.request_id, l.changed_by, u.full_name AS changer_name,
			l.from_status, l.to_status, l.note, l.created_at
		FROM request_status_logs l
		JOIN users u ON l.changed_by = u.id
		WHERE l.request_id = $1
		ORDER BY l.created_at ASC
	`
	rows, err := r.db.Query(query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []StatusLog
	for rows.Next() {
		var l StatusLog
		if err := rows.Scan(&l.ID, &l.RequestID, &l.ChangedBy, &l.ChangerName, &l.FromStatus, &l.ToStatus, &l.Note, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	if logs == nil {
		logs = []StatusLog{}
	}

	return logs, nil
}

// Create inserts a new maintenance request and records the initial audit log
func (r *Repository) Create(req MaintenanceRequest) (*MaintenanceRequest, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO maintenance_requests (asset_id, problem_description, priority, status, created_by)
		VALUES ($1, $2, $3, 'Submitted', $4)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(query, req.AssetID, req.ProblemDescription, req.Priority, req.CreatedBy).
		Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert request: %w", err)
	}
	req.Status = "Submitted"

	// Record initial creation in audit trail
	note := "Request created"
	logQuery := `
		INSERT INTO request_status_logs (request_id, changed_by, from_status, to_status, note)
		VALUES ($1, $2, NULL, 'Submitted', $3)
	`
	if _, err := tx.Exec(logQuery, req.ID, req.CreatedBy, note); err != nil {
		return nil, fmt.Errorf("failed to record initial audit log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &req, nil
}

// Update saves edits to a request
func (r *Repository) Update(id string, assetID, problemDescription, priority, status string) (*MaintenanceRequest, error) {
	query := `
		UPDATE maintenance_requests
		SET asset_id = $1, problem_description = $2, priority = $3, status = $4
		WHERE id = $5
		RETURNING id, asset_id, problem_description, priority, status, created_by, reviewed_by, reviewed_at, review_note, created_at, updated_at
	`
	var req MaintenanceRequest
	err := r.db.QueryRow(query, assetID, problemDescription, priority, status, id).
		Scan(&req.ID, &req.AssetID, &req.ProblemDescription, &req.Priority, &req.Status,
			&req.CreatedBy, &req.ReviewedBy, &req.ReviewedAt, &req.ReviewNote, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// Review updates request status (Approved/Rejected), reviewer ID, note, and adds an audit log
func (r *Repository) Review(id, reviewerID, newStatus, note, oldStatus string) (*MaintenanceRequest, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	query := `
		UPDATE maintenance_requests
		SET status = $1, reviewed_by = $2, reviewed_at = $3, review_note = $4
		WHERE id = $5
		RETURNING id, asset_id, problem_description, priority, status, created_by, reviewed_by, reviewed_at, review_note, created_at, updated_at
	`
	var req MaintenanceRequest
	err = tx.QueryRow(query, newStatus, reviewerID, now, note, id).
		Scan(&req.ID, &req.AssetID, &req.ProblemDescription, &req.Priority, &req.Status,
			&req.CreatedBy, &req.ReviewedBy, &req.ReviewedAt, &req.ReviewNote, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update review status: %w", err)
	}

	// Insert audit trail entry
	logQuery := `
		INSERT INTO request_status_logs (request_id, changed_by, from_status, to_status, note)
		VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := tx.Exec(logQuery, id, reviewerID, oldStatus, newStatus, note); err != nil {
		return nil, fmt.Errorf("failed to insert review audit log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &req, nil
}

// Delete removes a request
func (r *Repository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM maintenance_requests WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("request not found")
	}
	return nil
}

