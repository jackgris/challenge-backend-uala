CREATE TABLE IF NOT EXISTS likes (
    id TEXT PRIMARY KEY,         -- Auto-incrementing ID for the like entry
    tweet_id TEXT NOT NULL,        -- Associated tweet ID
    user_id TEXT NOT NULL,         -- User who liked
    FOREIGN KEY (tweet_id) REFERENCES tweets (id) ON DELETE CASCADE
);
