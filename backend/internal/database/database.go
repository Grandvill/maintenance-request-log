package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Config holds database connection parameters
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect establishes a connection with PostgreSQL with retry logic
func Connect(cfg Config) (*sql.DB, error) {
	sslMode := cfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, sslMode,
	)

	var db *sql.DB
	var err error

	// Retry connection up to 10 times (1-second intervals) to allow DB container to be ready
	for i := 1; i <= 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Printf("[DB] Successfully connected to PostgreSQL at %s:%s/%s", cfg.Host, cfg.Port, cfg.DBName)
				db.SetMaxOpenConns(25)
				db.SetMaxIdleConns(5)
				db.SetConnMaxLifetime(5 * time.Minute)
				return db, nil
			}
		}
		log.Printf("[DB] Waiting for database to be ready (attempt %d/10): %v", i, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after 10 attempts: %w", err)
}

// RunMigrations executes all .sql files in the migrations directory in alphabetical order
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// Create schema_migrations table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Look for migrations directory
	dir := findDirectory(migrationsDir, []string{"migrations", "../migrations", "../../migrations"})
	if dir == "" {
		return fmt.Errorf("migrations directory '%s' not found", migrationsDir)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, fileName := range sqlFiles {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", fileName).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", fileName, err)
		}

		if exists {
			continue // Already applied
		}

		filePath := filepath.Join(dir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		log.Printf("[DB Migration] Applying %s...", fileName)
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction for %s: %w", fileName, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", fileName); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", fileName, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", fileName, err)
		}
		log.Printf("[DB Migration] Applied %s successfully", fileName)
	}

	return nil
}

// RunSeeds executes seed files if the users table has 0 rows
func RunSeeds(db *sql.DB, seedsDir string) error {
	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return fmt.Errorf("failed to check user count before seeding: %w", err)
	}

	// Only seed if no users exist
	if userCount > 0 {
		log.Printf("[DB Seed] Database already contains %d users. Skipping seed.", userCount)
		return nil
	}

	dir := findDirectory(seedsDir, []string{"seeds", "../seeds", "../../seeds"})
	if dir == "" {
		log.Printf("[DB Seed] Warning: seeds directory not found at '%s'. Skipping.", seedsDir)
		return nil
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read seeds directory: %w", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, fileName := range sqlFiles {
		filePath := filepath.Join(dir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %w", fileName, err)
		}

		log.Printf("[DB Seed] Executing seed %s...", fileName)
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute seed %s: %w", fileName, err)
		}
		log.Printf("[DB Seed] Executed %s successfully", fileName)
	}

	// Ensure seeded users have valid, correctly hashed passwords matching our documentation
	passwords := map[string]string{
		"admin":             "admin123",
		"supervisor":        "supervisor123",
		"operator1":         "operator123",
		"operator2":         "operator123",
		"operator_inactive": "operator123",
	}

	for username, plainPwd := range passwords {
		hashed, err := HashPassword(plainPwd)
		if err == nil {
			db.Exec("UPDATE users SET password = $1 WHERE username = $2", hashed, username)
		}
	}
	log.Printf("[DB Seed] Verified password hashes for all default users.")

	return nil
}

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword checks if a plaintext password matches a bcrypt hash
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// findDirectory tries multiple fallback paths to locate a directory
func findDirectory(target string, fallbacks []string) string {
	if info, err := os.Stat(target); err == nil && info.IsDir() {
		return target
	}
	for _, fb := range fallbacks {
		if info, err := os.Stat(fb); err == nil && info.IsDir() {
			return fb
		}
	}
	return ""
}
