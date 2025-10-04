-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    daily_streak INT DEFAULT 0,
    total_checkins INT DEFAULT 0,
    total_messages INT DEFAULT 0,
    total_upvotes INT DEFAULT 0,
    total_upvotes_received INT DEFAULT 0,
    last_checkin_date DATE,
    last_upvote_given_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_stats;
-- +goose StatementEnd
