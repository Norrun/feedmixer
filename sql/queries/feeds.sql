-- name: AddFeed :one
INSERT INTO feeds ( 
    created_at, updated_at, name, url
    ) VALUES (
    current_timestamp,
    current_timestamp,
    ?,
    ?
)
RETURNING *;

-- name: GetFeedsByTag :many
SELECT * FROM feeds
WHERE id IN (
    SELECT feed_id FROM tags_feeds
    WHERE tag_id = ?
);

-- name: GetAllFeeds :many
SELECT * FROM feeds;

