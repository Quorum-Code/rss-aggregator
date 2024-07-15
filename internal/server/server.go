package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/Quorum-Code/rss-aggregator/internal/database"
	"github.com/Quorum-Code/rss-aggregator/internal/endpoints"
)

func StartServer(svrcfg ServerConfig) {
	// Get database conn
	dbURL := os.Getenv("connstr")
	if dbURL == "" {
		log.Fatal(errors.New("empty connstr"))
		return
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	dbQueries := database.New(db)

	// Build server
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    "localhost:" + svrcfg.Port,
	}
	apicfg := endpoints.ApiConfig{
		DB: dbQueries,
	}

	// Add endpoints
	mux.HandleFunc("GET /v1/healthz", apicfg.GetHealthz)
	mux.HandleFunc("GET /v1/err", apicfg.GetErr)

	mux.HandleFunc("POST /v1/users", apicfg.PostUsers)
	mux.HandleFunc("GET /v1/users", apicfg.MiddlewareAuth(apicfg.GetUserByAPIKey))

	mux.HandleFunc("POST /v1/feeds", apicfg.MiddlewareAuth(apicfg.PostFeed))

	// Start server
	server.ListenAndServe()
}
