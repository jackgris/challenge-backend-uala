CREATE TABLE IF NOT EXISTS retweets (
    id TEXT PRIMARY KEY,         -- Auto-incrementing ID for the retweet entry
    tweet_id TEXT NOT NULL,        -- Associated tweet ID
    user_id TEXT NOT NULL,         -- User who retweeted
    FOREIGN KEY (tweet_id) REFERENCES tweets (id) ON DELETE CASCADE
);
