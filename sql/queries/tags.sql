-- name: GetRelatedTags :many
SELECT * FROM tags
WHERE id IN (
    SELECT tags_feeds.tag_id FROM tags_feeds
    WHERE feed_id IN (
        SELECT feed_id FROM tags_feeds
        WHERE tags_feeds.tag_id = ?
    )
);


-- name: AddTag :one
INSERT INTO tags (
    created_at, updated_at, name
) VALUES (
    current_timestamp,
    current_timestamp,
    ?
)
RETURNING *;

-- name: GetTagByName :one
SELECT * FROM tags
WHERE name = ?;

-- name: GetTagById :one
SELECT * FROM tags
WHERE id = ?;



---- name: AddTags :exec
--INSERT INTO tags (created_at, updated_at, name)
--SELECT current_timestamp, current_timestamp, e.value) FROM json_each(sqlc.slice('name') AS e ;

-- name: GetTagTree :many
WITH RECURSIVE tag_tree AS (
    SELECT id, NULL AS parent_id, name, 0 AS level
    FROM tags
    UNION ALL
    SELECT t.id, tag_tree.id, t.name, tag_tree.level + 1
    FROM tags t
    WHERE t.id IN (
        SELECT tags_feeds.tag_id FROM tags_feeds
        WHERE feed_id IN (
        SELECT feed_id FROM tags_feeds
        WHERE tags_feeds.tag_id = tag_tree.id
        )
    )
)
SELECT * FROM tag_tree
ORDER BY name DESC, level ASC;

