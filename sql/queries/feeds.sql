-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedByUserID :one
SELECT *
FROM feeds
WHERE user_id = ($1);

-- name: GetFeedByID :one
SELECT *
FROM feeds
WHERE id = ($1);

-- name: GetFeeds :many
SELECT id, name, url, created_at, updated_at, last_fetched_at
FROM feeds;

-- name: GetFeedsWithNullFetched :many
SELECT *
FROM feeds
WHERE last_fetched_at IS NULL;

-- name: GetFeedsWithOldFetched :many
SELECT *
FROM feeds
WHERE last_fetched_at = (
    SELECT last_fetched_at
    FROM feeds
    ORDER BY last_fetched_at ASC
    LIMIT 1
);

-- name: GetNFeedsWithOldFetched :many
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT ($1);

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ($1)
RETURNING *;