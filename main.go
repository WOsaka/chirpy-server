package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/WOsaka/chirpy-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	cfg := &apiConfig{
		db: database.New(db),
		platform: os.Getenv("PLATFORM"),
		jwtSecret: os.Getenv("JWT_SECRET"),
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		http.StripPrefix("/app",
			cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthCheckHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	mux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)
	mux.HandleFunc("POST /api/login", cfg.loginHandler)
	mux.HandleFunc("POST /api/refresh", cfg.refreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", cfg.revokeRefreshTokenHandler)
	mux.HandleFunc("PUT /api/users", cfg.updateCredentialsHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirpHandler)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.setChirpyRedHandler)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	fmt.Println("Server listening on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}
