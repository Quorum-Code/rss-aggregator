-- +goose Up
CREATE TABLE feed_follows(
    feed_id UUID REFERENCES feeds(id),
    user_id UUID REFERENCES users(id),
    PRIMARY KEY (feed_id, user_id)
);

-- +goose Down
DROP TABLE feed_follows;