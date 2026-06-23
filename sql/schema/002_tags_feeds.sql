-- +goose Up
CREATE TABLE tags_feeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at VARCHAR(30) NOT NULL,
    updated_at VARCHAR(30) NOT NULL,
    tag_id INTEGER NOT NULL,
    feed_id INTEGER NOT NULL,
    UNIQUE(tag_id, feed_id),
    CONSTRAINT fk_feed_id 
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
     ON DELETE CASCADE,
     CONSTRAINT fk_tag_id 
    FOREIGN KEY (tag_id) REFERENCES tags(id)
     ON DELETE CASCADE
);