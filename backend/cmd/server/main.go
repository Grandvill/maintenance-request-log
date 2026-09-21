package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"maintenance-request-log/backend/internal/auth"
	"maintenance-request-log/backend/internal/database"
	"maintenance-request-log/backend/internal/middleware"
	"maintenance-request-log/backend/internal/requests"
	"maintenance-request-log/backend/internal/users"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("[Server] Starting Maintenance Request Log API Server...")

	// 1. Load Configurations from Environment
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbUser := getEnv("POSTGRES_USER", "postgres")
	dbPassword := getEnv("POSTGRES_PASSWORD", "postgres_secret")
	dbName := getEnv("POSTGRES_DB", "maintenance_db")

	port := getEnv("BACKEND_PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "hiroseelectric_super_secret_jwt_key_2026")
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	expiresHours, err := strconv.Atoi(getEnv("JWT_EXPIRES_IN_HOURS", "24"))
	if err != nil {
		expiresHours = 24
	}

	// 2. Connect to Database with Retry Logic
	dbConfig := database.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  "disable",
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("[FATAL] Could not connect to database: %v", err)
	}
	defer db.Close()

	// 3. Run Migrations Automatically
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("[FATAL] Migration failed: %v", err)
	}

	// 4. Run Seeds Automatically (only if DB is empty)
	if err := database.RunSeeds(db, "seeds"); err != nil {
		log.Printf("[WARN] Seed execution encountered an error: %v", err)
	}

	// 5. Initialize Layers
	// Auth
	authService := auth.NewService(db, jwtSecret, expiresHours)
	authHandler := auth.NewHandler(authService)

	// Users
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	// Requests
	requestRepo := requests.NewRepository(db)
	requestService := requests.NewService(requestRepo)
	requestHandler := requests.NewHandler(requestService)

	// 6. Router Setup
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.SetupCORS(corsOrigins))

	// Health Check Endpoint (Bonus Task: Health check endpoint)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		dbStatus := "connected"
		if err := db.Ping(); err != nil {
			status = "degraded"
			dbStatus = "disconnected"
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    status,
			"database":  dbStatus,
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		})
	})

	// API v1 Routes
	r.Route("/api/v1", func(api chi.Router) {
		// Public Auth Endpoints
		api.Post("/auth/login", authHandler.Login)
		api.Post("/auth/logout", authHandler.Logout)

		// Protected Routes (Require valid JWT)
		api.Group(func(protected chi.Router) {
			protected.Use(middleware.AuthRequired(jwtSecret))

			// Current User Profile
			protected.Get("/auth/me", authHandler.Me)

			// Maintenance Requests (Accessible by all roles; filtered internally per RBAC)
			protected.Route("/requests", func(reqRouter chi.Router) {
				reqRouter.Get("/", requestHandler.List)
				reqRouter.Post("/", requestHandler.Create)
				reqRouter.Get("/{id}", requestHandler.GetByID)
				reqRouter.Put("/{id}", requestHandler.Update)

				// Review: Supervisor and Admin ONLY
				reqRouter.With(middleware.RequireRoles("supervisor", "admin")).
					Patch("/{id}/review", requestHandler.Review)

				// Delete: Admin ONLY
				reqRouter.With(middleware.RequireRoles("admin")).
					Delete("/{id}", requestHandler.Delete)
			})

			// User Management (Admin ONLY)
			protected.Route("/users", func(userRouter chi.Router) {
				userRouter.Use(middleware.RequireRoles("admin"))

				userRouter.Get("/", userHandler.List)
				userRouter.Post("/", userHandler.Create)
				userRouter.Get("/{id}", userHandler.GetByID)
				userRouter.Put("/{id}", userHandler.Update)
				userRouter.Patch("/{id}/status", userHandler.ToggleStatus)
			})
		})
	})

	// 7. HTTP Server with Graceful Shutdown
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[Server] Listening on port %s (http://localhost:%s)", port, port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server error: %v", err)
		}
	}()

	// Wait for terminate signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Server] Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[FATAL] Server forced to shutdown: %v", err)
	}

	log.Println("[Server] Server exited cleanly.")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

