package requests

import "time"

// MaintenanceRequest entity matching the database schema
type MaintenanceRequest struct {
	ID                 string     `json:"id"`
	AssetID            string     `json:"asset_id"`
	ProblemDescription string     `json:"problem_description"`
	Priority           string     `json:"priority"`
	Status             string     `json:"status"`
	CreatedBy          string     `json:"created_by"`
	CreatorName        string     `json:"creator_name,omitempty"`
	ReviewedBy         *string    `json:"reviewed_by,omitempty"`
	ReviewerName       *string    `json:"reviewer_name,omitempty"`
	ReviewedAt         *time.Time `json:"reviewed_at,omitempty"`
	ReviewNote         *string    `json:"review_note,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// RequestDetail includes the request data and its audit trail logs
type RequestDetail struct {
	MaintenanceRequest
	AuditLogs []StatusLog `json:"audit_logs"`
}

// StatusLog represents an audit trail event
type StatusLog struct {
	ID          string    `json:"id"`
	RequestID   string    `json:"request_id"`
	ChangedBy   string    `json:"changed_by"`
	ChangerName string    `json:"changer_name"`
	FromStatus  *string   `json:"from_status"`
	ToStatus    string    `json:"to_status"`
	Note        *string   `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateRequestDTO input payload
type CreateRequestDTO struct {
	AssetID            string `json:"asset_id"`
	ProblemDescription string `json:"problem_description"`
	Priority           string `json:"priority"` // Low, Medium, High, Urgent
}

// UpdateRequestDTO input payload
type UpdateRequestDTO struct {
	AssetID            string  `json:"asset_id"`
	ProblemDescription string  `json:"problem_description"`
	Priority           string  `json:"priority"`
	Status             *string `json:"status,omitempty"` // Only Admin can directly edit status
}

// ReviewRequestDTO input payload for approve/reject
type ReviewRequestDTO struct {
	Status string `json:"status"` // Approved or Rejected
	Note   string `json:"note"`   // Optional reviewer note
}

// FilterParams for request list query
type FilterParams struct {
	Status    string
	Priority  string
	Search    string
	CreatedBy string // Used automatically when user is Operator
}

