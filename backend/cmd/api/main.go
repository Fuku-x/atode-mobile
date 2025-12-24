package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/atode/backend/internal/auth"
	"example.com/atode/backend/internal/db"
	"example.com/atode/backend/internal/handler"
	"example.com/atode/backend/internal/middleware"
	"example.com/atode/backend/internal/repository"
	"example.com/atode/backend/internal/service"
)

type config struct {
	addr string
}

func loadConfig() config {
	addr := os.Getenv("ATODE_API_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return config{addr: addr}
}

func main() {
	cfg := loadConfig()

	verifier, err := auth.NewFirebaseVerifier(context.Background())
	if err != nil {
		log.Fatalf("failed to init firebase verifier: %v", err)
	}

	dbConn, err := db.Open(context.Background())
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer func() {
		_ = dbConn.Close()
	}()

	// Run schema migrations on startup.
	// ATODE_MIGRATIONS_PATH can override the default "migrations" directory.
	if err := db.MigrateUp(context.Background(), os.Getenv("ATODE_DATABASE_URL")); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(dbConn)
	taskRepo := repository.NewTaskRepository(dbConn)
	taskSvc := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)
	taskItemHandler := handler.NewTaskItemHandler(taskSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		pingCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.HealthCheck(pingCtx, dbConn); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("unhealthy"))
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.Handle(
		"/me",
		middleware.RequireFirebaseAuth(verifier, userRepo)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, ok := auth.UserIDFromContext(r.Context())
				if !ok {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				uid, ok := auth.UIDFromContext(r.Context())
				if !ok {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				email, _ := auth.EmailFromContext(r.Context())

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				_ = json.NewEncoder(w).Encode(map[string]string{
					"user_id": userID.String(),
					"firebase_uid": string(uid),
					"email": email,
				})
			}),
		),
	)

	mux.Handle(
		"/tasks",
		middleware.RequireFirebaseAuth(verifier, userRepo)(
			taskHandler,
		),
	)

	mux.Handle(
		"/tasks/",
		middleware.RequireFirebaseAuth(verifier, userRepo)(
			taskItemHandler,
		),
	)

	srv := &http.Server{
		Addr:              cfg.addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("api listening on %s", cfg.addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
	case sig := <-sigCh:
		log.Printf("received signal: %s", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
