-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUserID :one
SELECT * FROM feeds
    WHERE user_id = $1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
    WHERE url = $1;

-- name: DeleteAllFeeds :exec
DELETE FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
    SET
        last_fetched_at = NOW(),
        updated_at = NOW()
    WHERE
        feeds.id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
    JOIN feed_follows ON feeds.id = feed_follows.feed_id
    WHERE feed_follows.user_id = $1
    ORDER BY last_fetched_at ASC NULLS FIRST
    LIMIT 1;