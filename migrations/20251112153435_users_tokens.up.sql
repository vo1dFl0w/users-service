CREATE TABLE users_tokens (
    user_id UUID NOT NULL,
    refresh_token TEXT NOT NULL,
    refresh_token_expiry TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);