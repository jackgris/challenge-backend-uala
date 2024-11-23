CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    follower_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0,
    salt TEXT,
    token TEXT,
    date_created TIMESTAMP NOT NULL,
    encoded_date TEXT NOT NULL
);
