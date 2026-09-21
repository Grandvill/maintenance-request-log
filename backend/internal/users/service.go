package users

import (
	"errors"
	"fmt"
	"strings"

	"maintenance-request-log/backend/internal/database"
)

var (
	ErrInvalidRole         = errors.New("invalid role: must be operator, supervisor, or admin")
	ErrUsernameRequired    = errors.New("username is required")
	ErrPasswordRequired    = errors.New("password must be at least 6 characters")
	ErrFullNameRequired    = errors.New("full name is required")
	ErrCannotDeactivateSelf = errors.New("you cannot deactivate your own admin account")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListUsers() ([]User, error) {
	return s.repo.ListAll()
}

func (s *Service) GetUser(id string) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) CreateUser(req CreateUserRequest) (*User, error) {
	if strings.TrimSpace(req.Username) == "" {
		return nil, ErrUsernameRequired
	}
	if strings.TrimSpace(req.FullName) == "" {
		return nil, ErrFullNameRequired
	}
	if len(req.Password) < 6 {
		return nil, ErrPasswordRequired
	}

	role := strings.ToLower(strings.TrimSpace(req.Role))
	if role != "operator" && role != "supervisor" && role != "admin" {
		return nil, ErrInvalidRole
	}

	hashedPassword, err := database.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := User{
		Username: strings.TrimSpace(req.Username),
		Password: hashedPassword,
		FullName: strings.TrimSpace(req.FullName),
		Role:     role,
		IsActive: true,
	}

	return s.repo.Create(u)
}

func (s *Service) UpdateUser(id string, req UpdateUserRequest) (*User, error) {
	if strings.TrimSpace(req.FullName) == "" {
		return nil, ErrFullNameRequired
	}

	role := strings.ToLower(strings.TrimSpace(req.Role))
	if role != "operator" && role != "supervisor" && role != "admin" {
		return nil, ErrInvalidRole
	}

	var hashedPassword *string
	if req.Password != nil && strings.TrimSpace(*req.Password) != "" {
		if len(*req.Password) < 6 {
			return nil, ErrPasswordRequired
		}
		hash, err := database.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashedPassword = &hash
	}

	return s.repo.Update(id, strings.TrimSpace(req.FullName), role, hashedPassword)
}

func (s *Service) ToggleStatus(targetUserID string, currentUserID string, isActive bool) (*User, error) {
	if targetUserID == currentUserID && !isActive {
		return nil, ErrCannotDeactivateSelf
	}
	return s.repo.SetStatus(targetUserID, isActive)
}

