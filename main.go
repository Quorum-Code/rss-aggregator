package main

import (
	"fmt"
	"os"

	"github.com/Quorum-Code/rss-aggregator/internal/server"
	"github.com/joho/godotenv"
)

var EnvFile = "./.env"

func main() {
	// Load .env file
	godotenv.Load(EnvFile)

	port := os.Getenv("PORT")

	fmt.Printf("Serving at localhost:%s\n", port)

	svrcfg := server.ServerConfig{Port: port}
	server.StartServer(svrcfg)
}
