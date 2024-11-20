CREATE TABLE IF NOT EXISTS tweets (
    id TEXT PRIMARY KEY,           -- ID of the tweet
    user_id TEXT NOT NULL,        -- Username of the author
    content TEXT NOT NULL,         -- Content of the tweet
    created_at TIMESTAMP NOT NULL, -- Creation timestamp
    encoded_date TEXT NOT NULL,    -- Encoded date as a string
    like_count INT DEFAULT 0,      -- Number of likes
    retweet_count INT DEFAULT 0    -- Number of retweets
);
