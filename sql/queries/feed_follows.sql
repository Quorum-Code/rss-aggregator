-- name: CreateFeedFollow :one
INSERT INTO feed_follows (feed_id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetFeedFollows :many
SELECT *
FROM feed_follows;
