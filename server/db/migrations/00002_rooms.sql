-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    creator_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    is_pinned BOOLEAN DEFAULT FALSE,
    topic_title TEXT,
    topic_description TEXT,
    topic_url TEXT,
    topic_source TEXT,
    topic_updated_at TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rooms;
-- +goose StatementEnd