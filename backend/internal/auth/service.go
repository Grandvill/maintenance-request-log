package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"maintenance-request-log/backend/internal/database"
	"maintenance-request-log/backend/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrAccountDeactivated = errors.New("your account has been deactivated, please contact an administrator")
)

type Service struct {
	db           *sql.DB
	jwtSecret    string
	expiresHours int
}

func NewService(db *sql.DB, jwtSecret string, expiresHours int) *Service {
	if expiresHours <= 0 {
		expiresHours = 24
	}
	return &Service{
		db:           db,
		jwtSecret:    jwtSecret,
		expiresHours: expiresHours,
	}
}

// Authenticate verifies credentials and returns a signed JWT with user profile
func (s *Service) Authenticate(req LoginRequest) (*LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, ErrInvalidCredentials
	}

	query := `
		SELECT id, username, password, full_name, role, is_active, created_at
		FROM users
		WHERE username = $1
	`

	var (
		id, username, passwordHash, fullName, role string
		isActive                                   bool
		createdAt                                  time.Time
	)

	err := s.db.QueryRow(query, req.Username).Scan(
		&id, &username, &passwordHash, &fullName, &role, &isActive, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("database query error: %w", err)
	}

	// Check if account is active
	if !isActive {
		return nil, ErrAccountDeactivated
	}

	// Verify password
	if !database.CheckPassword(passwordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT
	expiresAt := time.Now().Add(time.Duration(s.expiresHours) * time.Hour)
	claims := middleware.CustomClaims{
		UserID:   id,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		User: UserResponse{
			ID:        id,
			Username:  username,
			FullName:  fullName,
			Role:      role,
			IsActive:  isActive,
			CreatedAt: createdAt,
		},
	}, nil
}

// GetUserByID fetches user profile by ID
func (s *Service) GetUserByID(userID string) (*UserResponse, error) {
	query := `
		SELECT id, username, full_name, role, is_active, created_at
		FROM users
		WHERE id = $1
	`
	var user UserResponse
	err := s.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.FullName, &user.Role, &user.IsActive, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

