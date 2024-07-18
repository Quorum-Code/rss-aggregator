package endpoints

import (
	"encoding/xml"
	"fmt"
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
	dbfeeds := cfg.getNextFeedsToFetch(req)
	feeds := DatabaseFeedsToFeeds(dbfeeds)

	// For each
	fmt.Println("Printing feed text")
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

		fmt.Println(rss.Channel.Title)
	}
	fmt.Println("Done printing feed text")

	// Read XML from url

	// Read/Parse XML from feed
	// var r database.RSS
	// err = xml.Unmarshal(nil, &r)
	// if err != nil {
	// 	respondWithError(resp, http.StatusInternalServerError, err)
	// 	return
	// }

	// Update last_fetched_at and updated_at

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

func (cfg *ApiConfig) getNextFeedsToFetch(req *http.Request) []database.Feed {
	feeds, err := cfg.DB.GetFeedsWithNullFetched(req.Context())
	if err == nil {
		return feeds
	} else {
		fmt.Println(err.Error())
	}

	feeds, err = cfg.DB.GetFeedsWithOldFetched(req.Context())
	if err == nil {
		return feeds
	}

	fmt.Println(err.Error())
	return nil
}
