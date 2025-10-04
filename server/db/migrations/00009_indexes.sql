-- +goose Up

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages(room_id);
CREATE INDEX IF NOT EXISTS idx_rooms_expires_at ON rooms(expires_at);
CREATE INDEX IF NOT EXISTS idx_upvotes_from_user ON upvotes(from_user_id);
CREATE INDEX IF NOT EXISTS idx_upvotes_to_user ON upvotes(to_user_id);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_room_id;
DROP INDEX IF EXISTS idx_rooms_expires_at;
DROP INDEX IF EXISTS idx_upvotes_from_user;
DROP INDEX IF EXISTS idx_upvotes_to_user;
-- +goose StatementEnd