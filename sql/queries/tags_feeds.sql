-- name: GetTagAndFeedIds :many
SELECT * FROM tags_feeds
ORDER BY feed_id;

-- name: AttachTag :one
INSERT INTO tags_feeds ( 
    created_at, updated_at, feed_id, tag_id
    ) VALUES (
    current_timestamp,
    current_timestamp,
    ?,
    ?
)
RETURNING *;