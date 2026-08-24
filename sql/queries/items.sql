-- name: GetAllItems :many
SELECT * FROM items;

-- name: AddItem :one
INSERT INTO items ( 
    created_at, updated_at, title, url, description, published_at,feed_id
    ) VALUES (
    current_timestamp,
    current_timestamp,
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;