CREATE TABLE users_scoreboard (
    user_id UUID PRIMARY KEY REFERENCES users (user_id) ON DELETE CASCADE,
    score BIGINT NOT NULL DEFAULT 0 CHECK (score >= 0)
);