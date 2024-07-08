package server

import (
	"net/http"

	"github.com/Quorum-Code/rss-aggregator/internal/endpoints"
)

func StartServer(svrcfg ServerConfig) {
	// Build server
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    "localhost:" + svrcfg.Port,
	}
	apicfg := endpoints.ApiConfig{}

	// Add endpoints
	mux.HandleFunc("GET /v1/healthz", apicfg.GetHealthz)
	mux.HandleFunc("GET /v1/err", apicfg.GetErr)

	// Start server
	server.ListenAndServe()
}
