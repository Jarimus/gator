-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
  )
  RETURNING *
)

SELECT
    inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
FROM
    inserted_feed_follow
INNER JOIN feeds
    ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users
    ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUserByID :many
SELECT
    feed_follows.*, users.name AS user_name, feeds.name AS feed_name, feeds.url as url
FROM
    feed_follows
INNER JOIN feeds
    ON feed_follows.feed_id = feeds.id
INNER JOIN users
    ON feed_follows.user_id = users.id
WHERE
    feed_follows.user_id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
    WHERE
        feed_follows.feed_id = $1 AND
        feed_follows.user_id = $2;