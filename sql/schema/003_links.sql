-- +goose Up
CREATE TABLE links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at VARCHAR(30) NOT NULL,
    updated_at VARCHAR(30) NOT NULL,
    url TEXT NOT NULL,
    item_id INTEGER,
    CONSTRAINT fk_item_id
    FOREIGN KEY (item_id) REFERENCES items(id)
     ON DELETE CASCADE
);