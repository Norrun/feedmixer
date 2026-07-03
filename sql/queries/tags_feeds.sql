-- name: GetTagAndFeedIds :many
SELECT * FROM tags_feeds
ORDER BY feed_id;