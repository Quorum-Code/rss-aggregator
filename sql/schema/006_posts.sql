-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY,
    feed_id UUID REFERENCES feeds(id),
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE posts;