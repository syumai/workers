-- Migration number: 0000 	 2023-01-09T14:48:53.705Z
CREATE TABLE articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at INTEGER NOT NULL
);
CREATE INDEX idx_articles_on_created_at ON articles (created_at DESC);
