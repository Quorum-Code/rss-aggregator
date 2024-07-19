-- name: CreatePost :one
INSERT INTO posts (id, feed_id, title, url, description)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPostByID :one
SELECT *
FROM posts
WHERE id = ($1);

-- name: GetPostsByUser :many
SELECT *
FROM posts
WHERE feed_id IN (
    SELECT feed_id
    FROM feed_follows
    WHERE user_id = ($1) 
);