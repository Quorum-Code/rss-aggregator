-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY,
    feed_id UUID REFERENCES feeds(id),
    title TEXT,
    url TEXT,
    description TEXT,
    published_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE posts;