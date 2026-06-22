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