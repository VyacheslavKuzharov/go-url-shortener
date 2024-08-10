CREATE TABLE IF NOT EXISTS shorten_urls(
    id serial PRIMARY KEY,
    short_key VARCHAR (255) UNIQUE NOT NULL,
    original_url VARCHAR (255) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_short_key ON shorten_urls(short_key);