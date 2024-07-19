-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY,
    feed_id UUID REFERENCES feeds(id) NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT(CURRENT_TIMESTAMP) NOT NULL,
    updated_at TIMESTAMP DEFAULT(CURRENT_TIMESTAMP) NOT NULL
);

-- +goose Down
DROP TABLE posts;