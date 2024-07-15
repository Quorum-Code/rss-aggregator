-- name: CreateFeedFollow :one
INSERT INTO feed_follows (feed_id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetFeedFollows :many
SELECT *
FROM feed_follows;

-- name: GetFeedFollowsByUserID :many
SELECT *
FROM feed_follows
WHERE user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE feed_id = $1 AND user_id = $2;