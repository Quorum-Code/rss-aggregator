package endpoints

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Quorum-Code/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) GetFeeds(resp http.ResponseWriter, req *http.Request) {
	feeds, err := cfg.DB.GetFeeds(req.Context())
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(resp, http.StatusAccepted, feeds)
}

func (cfg *ApiConfig) PostFeed(resp http.ResponseWriter, req *http.Request, user database.User) {
	type body struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	var b body
	err := jsonUnmarshalReader(req.Body, &b)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	// Create a feed entry
	feed, err := cfg.DB.CreateFeed(req.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      b.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Url:       b.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	// Create a follow entry
	follow, err := cfg.DB.CreateFeedFollow(req.Context(), database.CreateFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	// Wrap feed and follow into interface
	type response struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}
	r := response{Feed: feed, FeedFollow: follow}
	respondWithJSON(resp, http.StatusOK, r)
}

func (cfg *ApiConfig) RefreshFetches(resp http.ResponseWriter, req *http.Request) {
	// Get feeds to refresh
	dbfeeds, err := cfg.getNextFeedsToFetch(req)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	feeds := DatabaseFeedsToFeeds(dbfeeds)

	// For each
	fmt.Printf(" * Refreshing Feeds [%d] * \n", len(feeds))
	for i := 0; i < len(feeds); i++ {
		res, err := http.Get(feeds[i].Url)
		if err != nil {
			respondWithError(resp, http.StatusInternalServerError, err)
			return
		}

		// Unmarshal
		rss := database.RSS{}
		decoder := xml.NewDecoder(res.Body)
		err = decoder.Decode(&rss)
		if err != nil {
			respondWithError(resp, http.StatusInternalServerError, err)
			return
		}

		dbfeed, err := cfg.DB.MarkFeedFetched(req.Context(), feeds[i].ID)
		if err != nil {
			respondWithError(resp, http.StatusInternalServerError, err)
			return
		}
		feeds[i] = DatabaseFeedToFeed(dbfeed)

		fmt.Printf(" - %s\n", rss.Channel.Title)
	}

	// Respond with feeds
	respondWithJSON(resp, http.StatusOK, feeds)
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func DatabaseFeedsToFeeds(dbfeeds []database.Feed) []Feed {
	if dbfeeds == nil {
		return make([]Feed, 0)
	}

	feeds := make([]Feed, 0, len(dbfeeds))

	for _, dbf := range dbfeeds {
		feeds = append(feeds, DatabaseFeedToFeed(dbf))
	}

	return feeds
}

func DatabaseFeedToFeed(feed database.Feed) Feed {
	var lastfetchedat *time.Time
	if feed.LastFetchedAt.Valid {
		lastfetchedat = &feed.LastFetchedAt.Time
	}

	return Feed{
		ID:            feed.ID,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		LastFetchedAt: lastfetchedat,
	}
}

func (cfg *ApiConfig) getNextFeedsToFetch(req *http.Request) ([]database.Feed, error) {
	// Get not yet fetched
	feeds, err := cfg.DB.GetFeedsWithNullFetched(req.Context())
	if len(feeds) > 0 {
		return feeds, nil
	} else if err != nil {
		log.Fatal(err)
	}

	// Get oldest fetched
	feeds, err = cfg.DB.GetFeedsWithOldFetched(req.Context())
	if len(feeds) > 0 {
		return feeds, nil
	} else if err != nil {
		log.Fatal(err)
	}

	return nil, errors.New("no feeds to fetch")
}
