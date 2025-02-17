// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: posts.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (id, feed_id, title, url, description)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, feed_id, title, url, description, published_at, created_at, updated_at
`

type CreatePostParams struct {
	ID          uuid.UUID
	FeedID      uuid.UUID
	Title       string
	Url         string
	Description string
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.FeedID,
		arg.Title,
		arg.Url,
		arg.Description,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.FeedID,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPostByID = `-- name: GetPostByID :one
SELECT id, feed_id, title, url, description, published_at, created_at, updated_at
FROM posts
WHERE id = ($1)
`

func (q *Queries) GetPostByID(ctx context.Context, id uuid.UUID) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostByID, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.FeedID,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPostsByUser = `-- name: GetPostsByUser :many
SELECT id, feed_id, title, url, description, published_at, created_at, updated_at
FROM posts
WHERE feed_id IN (
    SELECT feed_id
    FROM feed_follows
    WHERE user_id = ($1) 
)
`

func (q *Queries) GetPostsByUser(ctx context.Context, userID uuid.UUID) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.FeedID,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
