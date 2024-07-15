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

	// Users
	mux.HandleFunc("POST /v1/users", apicfg.PostUsers)
	mux.HandleFunc("GET /v1/users", apicfg.MiddlewareAuth(apicfg.GetUserByAPIKey))

	// Feeds
	mux.HandleFunc("POST /v1/feeds", apicfg.MiddlewareAuth(apicfg.PostFeed))
	mux.HandleFunc("GET /v1/feeds", apicfg.GetFeeds)

	// Feed follows
	mux.HandleFunc("GET /v1/feed_follows", apicfg.MiddlewareAuth(apicfg.GetFeedFollows))
	mux.HandleFunc("POST /v1/feed_follows", apicfg.MiddlewareAuth(apicfg.PostFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apicfg.MiddlewareAuth(apicfg.DeleteFeedFollow))

	mux.HandleFunc("GET /v1/refresh_feeds", apicfg.RefreshFetches)

	// Start server
	server.ListenAndServe()
}
