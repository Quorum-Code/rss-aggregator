package endpoints

import "github.com/Quorum-Code/rss-aggregator/internal/database"

type ApiConfig struct {
	DB *database.Queries
}
