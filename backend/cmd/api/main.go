package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pfrancistechs-cel/algosearch/backend/internal/db"
	"github.com/pfrancistechs-cel/algosearch/backend/internal/search"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://algosearch:algosearch_password@localhost:5433/algosearch?sslmode=disable"
	}

	log.Printf("[AlgoSearch API] Booting on port %s...\n", port)

	// Connect to PostgreSQL (Neon compatible) and run migrations + seeding
	pool, err := db.ConnectAndMigrate(databaseURL)
	if err != nil {
		log.Fatalf("[AlgoSearch API] Fatal database error: %v\n", err)
	}
	defer pool.Close()

	// Initialize dependencies
	searchRepo := search.NewRepository(pool)
	searchService := search.NewService(searchRepo)
	searchHandler := search.NewHandler(searchService)

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","engine":"postgres-pg_trgm-tsvector","service":"algosearch"}`)
	})

	// Search and Document Endpoints
	searchHandler.RegisterRoutes(mux)

	// CORS and Logging Middleware
	corsHandler := withCORS(mux)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server runner
	go func() {
		log.Printf("[AlgoSearch API] SearchLite engine listening on http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[AlgoSearch API] ListenAndServe failed: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("[AlgoSearch API] Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[AlgoSearch API] Forced shutdown error: %v\n", err)
	}
	log.Println("[AlgoSearch API] Server exited cleanly.")
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
