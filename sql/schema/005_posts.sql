-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    title TEXT,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;