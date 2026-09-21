package requests

import (
	"testing"

	"maintenance-request-log/backend/internal/middleware"
)

// TestRBACMatrix verifies that the permission matrix is strictly enforced on the server
func TestRBACMatrix(t *testing.T) {
	// Dummy repository is not called for validation failures that happen before DB queries
	service := NewService(nil)

	operatorUser := middleware.AuthUser{
		ID:       "op-1",
		Username: "operator1",
		Role:     "operator",
	}

	supervisorUser := middleware.AuthUser{
		ID:       "sup-1",
		Username: "supervisor",
		Role:     "supervisor",
	}

	// 1. Operator attempts to review (Approve/Reject) -> MUST be Forbidden
	t.Run("Operator cannot review requests", func(t *testing.T) {
		_, err := service.ReviewRequest(operatorUser, "req-1", ReviewRequestDTO{
			Status: "Approved",
			Note:   "Operator trying to approve",
		})
		if err != ErrForbidden {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	// 2. Operator attempts to delete a request -> MUST be Forbidden
	t.Run("Operator cannot delete requests", func(t *testing.T) {
		err := service.DeleteRequest(operatorUser, "req-1")
		if err != ErrForbidden {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	// 3. Supervisor attempts to delete a request -> MUST be Forbidden
	t.Run("Supervisor cannot delete requests", func(t *testing.T) {
		err := service.DeleteRequest(supervisorUser, "req-1")
		if err != ErrForbidden {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	// 4. Invalid review status -> MUST return ErrInvalidReviewStatus
	t.Run("Review with invalid status is rejected", func(t *testing.T) {
		_, err := service.ReviewRequest(supervisorUser, "req-1", ReviewRequestDTO{
			Status: "PendingReview", // Invalid status
		})
		if err != ErrInvalidReviewStatus {
			t.Fatalf("expected ErrInvalidReviewStatus, got: %v", err)
		}
	})

	// 5. Validation helper tests
	t.Run("Priority validation", func(t *testing.T) {
		validPriorities := []string{"Low", "Medium", "High", "Urgent"}
		for _, p := range validPriorities {
			if !isValidPriority(p) {
				t.Errorf("expected %s to be valid priority", p)
			}
		}

		invalidPriorities := []string{"Critical", "Immediate", "Normal", "unknown", ""}
		for _, p := range invalidPriorities {
			if isValidPriority(p) {
				t.Errorf("expected %s to be invalid priority", p)
			}
		}
	})
}

