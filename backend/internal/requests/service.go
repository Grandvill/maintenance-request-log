package requests

import (
	"errors"
	"fmt"
	"strings"

	"maintenance-request-log/backend/internal/middleware"
)

var (
	ErrNotFound             = errors.New("maintenance request not found")
	ErrForbidden            = errors.New("forbidden: insufficient permissions")
	ErrAssetIDRequired      = errors.New("asset ID is required")
	ErrProblemDescRequired  = errors.New("problem description is required")
	ErrInvalidPriority      = errors.New("priority must be Low, Medium, High, or Urgent")
	ErrInvalidReviewStatus  = errors.New("review status must be Approved or Rejected")
	ErrOnlySubmittedEditable = errors.New("cannot edit request: only requests with 'Submitted' status can be modified")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ListRequests retrieves requests enforced by caller's role
func (s *Service) ListRequests(user middleware.AuthUser, params FilterParams) ([]MaintenanceRequest, error) {
	// Rule: Operator can ONLY view their own requests (enforced server-side)
	if user.Role == "operator" {
		params.CreatedBy = user.ID
	}
	// Supervisors and Admins can see all requests (or filter as provided)

	return s.repo.List(params)
}

// GetRequestDetail retrieves request with audit logs, strictly enforcing RBAC
func (s *Service) GetRequestDetail(user middleware.AuthUser, id string) (*RequestDetail, error) {
	req, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	// Rule: Operator cannot view other users' requests
	if user.Role == "operator" && req.CreatedBy != user.ID {
		return nil, ErrForbidden
	}

	logs, err := s.repo.GetAuditLogs(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch audit logs: %w", err)
	}

	return &RequestDetail{
		MaintenanceRequest: *req,
		AuditLogs:          logs,
	}, nil
}

// CreateRequest handles new request submission
func (s *Service) CreateRequest(user middleware.AuthUser, dto CreateRequestDTO) (*MaintenanceRequest, error) {
	if strings.TrimSpace(dto.AssetID) == "" {
		return nil, ErrAssetIDRequired
	}
	if strings.TrimSpace(dto.ProblemDescription) == "" {
		return nil, ErrProblemDescRequired
	}

	priority := capitalize(strings.TrimSpace(dto.Priority))
	if priority == "" {
		priority = "Medium"
	}
	if !isValidPriority(priority) {
		return nil, ErrInvalidPriority
	}

	req := MaintenanceRequest{
		AssetID:            strings.TrimSpace(dto.AssetID),
		ProblemDescription: strings.TrimSpace(dto.ProblemDescription),
		Priority:           priority,
		CreatedBy:          user.ID,
	}

	return s.repo.Create(req)
}

// UpdateRequest enforces rules:
// - Operator / Supervisor: can only edit OWN request while status is 'Submitted'
// - Admin: can edit ANY request in any status
func (s *Service) UpdateRequest(user middleware.AuthUser, id string, dto UpdateRequestDTO) (*MaintenanceRequest, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	// Permission checks
	if user.Role == "admin" {
		// Admin can edit any request
	} else {
		// Operator & Supervisor can only edit their OWN requests
		if existing.CreatedBy != user.ID {
			return nil, ErrForbidden
		}
		// Must be in Submitted status
		if existing.Status != "Submitted" {
			return nil, ErrOnlySubmittedEditable
		}
	}

	assetID := existing.AssetID
	if strings.TrimSpace(dto.AssetID) != "" {
		assetID = strings.TrimSpace(dto.AssetID)
	}

	desc := existing.ProblemDescription
	if strings.TrimSpace(dto.ProblemDescription) != "" {
		desc = strings.TrimSpace(dto.ProblemDescription)
	}

	priority := existing.Priority
	if strings.TrimSpace(dto.Priority) != "" {
		newPri := capitalize(strings.TrimSpace(dto.Priority))
		if !isValidPriority(newPri) {
			return nil, ErrInvalidPriority
		}
		priority = newPri
	}

	// Status can only be directly edited by Admin
	status := existing.Status
	if user.Role == "admin" && dto.Status != nil && strings.TrimSpace(*dto.Status) != "" {
		st := capitalize(strings.TrimSpace(*dto.Status))
		if st == "Submitted" || st == "Approved" || st == "Rejected" {
			status = st
		}
	}

	return s.repo.Update(id, assetID, desc, priority, status)
}

// ReviewRequest approves or rejects a request
// - Operator: Forbidden
// - Supervisor / Admin: Allowed
func (s *Service) ReviewRequest(user middleware.AuthUser, id string, dto ReviewRequestDTO) (*MaintenanceRequest, error) {
	if user.Role != "supervisor" && user.Role != "admin" {
		return nil, ErrForbidden
	}

	newStatus := capitalize(strings.TrimSpace(dto.Status))
	if newStatus != "Approved" && newStatus != "Rejected" {
		return nil, ErrInvalidReviewStatus
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	return s.repo.Review(id, user.ID, newStatus, strings.TrimSpace(dto.Note), existing.Status)
}

// DeleteRequest deletes a request
// - Admin ONLY
func (s *Service) DeleteRequest(user middleware.AuthUser, id string) error {
	if user.Role != "admin" {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

func isValidPriority(p string) bool {
	switch p {
	case "Low", "Medium", "High", "Urgent":
		return true
	default:
		return false
	}
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	lower := strings.ToLower(s)
	return strings.ToUpper(lower[:1]) + lower[1:]
}

