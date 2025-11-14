CREATE TABLE users_tasks (
    user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    task TEXT,
    reward BIGINT NOT NULL DEFAULT 0 CHECK (reward >= 0),
    complete BOOLEAN DEFAULT false,
    referrer_id UUID NULL REFERENCES users (user_id) ON DELETE SET NULL,
    PRIMARY KEY (user_id, task)
);

CREATE INDEX idx_users_tasks_user_id ON users_tasks (user_id);