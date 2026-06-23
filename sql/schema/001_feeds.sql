-- +goose Up
-- +goose statementbegin
CREATE TABLE feeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at VARCHAR(30) NOT NULL,
    updated_at VARCHAR(30) NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    last_fetched_at VARCHAR(30),
    last_checked_at VARCHAR(30)
);
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at VARCHAR(30) NOT NULL,
    updated_at VARCHAR(30) NOT NULL,
    external_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT,
    published_at VARCHAR(30),
    feed_id INTEGER NOT NULL,
    constraint fk_feed_id 
    foreign key (feed_id) REFERENCES feeds(id),
    UNIQUE(external_id, feed_id)
);
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at VARCHAR(30) NOT NULL,
    updated_at VARCHAR(30) NOT NULL,
    name VARCHAR(40) UNIQUE NOT NULL,
    last_checked_at VARCHAR(30)

);

-- +goose statementend

-- +goose Down
DROP TABLE feeds;