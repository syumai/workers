DROP TABLE IF EXISTS articles;
CREATE TABLE IF NOT EXISTS articles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at INT NOT NULL
);
CREATE INDEX idx_articles_on_created_at ON articles (created_at DESC);
INSERT INTO articles (title, body, created_at) VALUES (
    'title of example post',
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
    UNIX_TIMESTAMP()
);