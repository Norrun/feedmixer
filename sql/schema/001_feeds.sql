-- +goose Up
CREATE TABLE feeds (
    id INT PRIMARY KEY,
    created_at VARCHAR(30) NOT NULL,
    updated_at VARCHAR(30) NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    last_fetched_at VARCHAR(30)
);

-- +goose Down
DROP TABLE feeds;