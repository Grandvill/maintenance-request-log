package users

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ListAll returns all users ordered by creation date
func (r *Repository) ListAll() ([]User, error) {
	query := `
		SELECT id, username, full_name, role, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var userList []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		userList = append(userList, u)
	}

	if userList == nil {
		userList = []User{}
	}

	return userList, nil
}

// GetByID fetches a single user by ID
func (r *Repository) GetByID(id string) (*User, error) {
	query := `
		SELECT id, username, full_name, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

// Create inserts a new user
func (r *Repository) Create(u User) (*User, error) {
	query := `
		INSERT INTO users (username, password, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, u.Username, u.Password, u.FullName, u.Role, u.IsActive).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	return &u, nil
}

// Update updates full_name, role, and optionally password
func (r *Repository) Update(id string, fullName, role string, passwordHash *string) (*User, error) {
	var query string
	var err error

	if passwordHash != nil && *passwordHash != "" {
		query = `
			UPDATE users
			SET full_name = $1, role = $2, password = $3
			WHERE id = $4
			RETURNING id, username, full_name, role, is_active, created_at, updated_at
		`
		var u User
		err = r.db.QueryRow(query, fullName, role, *passwordHash, id).
			Scan(&u.ID, &u.Username, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &u, nil
	}

	query = `
		UPDATE users
		SET full_name = $1, role = $2
		WHERE id = $3
		RETURNING id, username, full_name, role, is_active, created_at, updated_at
	`
	var u User
	err = r.db.QueryRow(query, fullName, role, id).
		Scan(&u.ID, &u.Username, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// SetStatus toggles is_active for a user
func (r *Repository) SetStatus(id string, isActive bool) (*User, error) {
	query := `
		UPDATE users
		SET is_active = $1
		WHERE id = $2
		RETURNING id, username, full_name, role, is_active, created_at, updated_at
	`
	var u User
	err := r.db.QueryRow(query, isActive, id).
		Scan(&u.ID, &u.Username, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

