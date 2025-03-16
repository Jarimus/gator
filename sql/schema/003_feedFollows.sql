-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    feed_id UUID NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT id_pairs UNIQUE (feed_id, user_id)

);

-- +goose Down
DROP TABLE feed_follows;