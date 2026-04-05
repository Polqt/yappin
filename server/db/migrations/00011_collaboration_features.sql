-- +goose Up

-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS room_members (
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member',
    can_manage_room BOOLEAN NOT NULL DEFAULT FALSE,
    can_manage_channels BOOLEAN NOT NULL DEFAULT FALSE,
    can_moderate BOOLEAN NOT NULL DEFAULT FALSE,
    can_post BOOLEAN NOT NULL DEFAULT TRUE,
    muted_until TIMESTAMP,
    banned_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (room_id, user_id)
);

CREATE TABLE IF NOT EXISTS room_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    position INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS room_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    category_id UUID REFERENCES room_categories(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    kind TEXT NOT NULL DEFAULT 'text',
    position INT NOT NULL DEFAULT 0,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (room_id, name)
);

ALTER TABLE messages
    ADD COLUMN IF NOT EXISTS channel_id UUID REFERENCES room_channels(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS parent_message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
    message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    kind TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO room_members (
    room_id,
    user_id,
    role,
    can_manage_room,
    can_manage_channels,
    can_moderate,
    can_post
)
SELECT r.id, r.creator_id, 'owner', TRUE, TRUE, TRUE, TRUE
FROM rooms r
WHERE r.creator_id IS NOT NULL
ON CONFLICT (room_id, user_id) DO NOTHING;

INSERT INTO room_categories (room_id, name, position)
SELECT r.id, 'General', 0
FROM rooms r
WHERE NOT EXISTS (
    SELECT 1
    FROM room_categories rc
    WHERE rc.room_id = r.id
);

INSERT INTO room_channels (room_id, category_id, name, description, kind, position)
SELECT r.id, rc.id, 'lobby', 'Default channel', 'text', 0
FROM rooms r
JOIN LATERAL (
    SELECT id
    FROM room_categories
    WHERE room_id = r.id
    ORDER BY position ASC, created_at ASC
    LIMIT 1
) rc ON TRUE
WHERE NOT EXISTS (
    SELECT 1
    FROM room_channels ch
    WHERE ch.room_id = r.id
);

UPDATE messages m
SET channel_id = ch.id
FROM LATERAL (
    SELECT id
    FROM room_channels
    WHERE room_id = m.room_id
    ORDER BY position ASC, created_at ASC
    LIMIT 1
) ch
WHERE m.channel_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_room_members_user_id ON room_members(user_id);
CREATE INDEX IF NOT EXISTS idx_room_categories_room_id ON room_categories(room_id);
CREATE INDEX IF NOT EXISTS idx_room_channels_room_id ON room_channels(room_id);
CREATE INDEX IF NOT EXISTS idx_messages_channel_id ON messages(channel_id);
CREATE INDEX IF NOT EXISTS idx_messages_parent_message_id ON messages(parent_message_id);
CREATE INDEX IF NOT EXISTS idx_messages_content_search ON messages USING gin (content gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id_read ON notifications(user_id, is_read, created_at DESC);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_notifications_user_id_read;
DROP INDEX IF EXISTS idx_messages_content_search;
DROP INDEX IF EXISTS idx_messages_parent_message_id;
DROP INDEX IF EXISTS idx_messages_channel_id;
DROP INDEX IF EXISTS idx_room_channels_room_id;
DROP INDEX IF EXISTS idx_room_categories_room_id;
DROP INDEX IF EXISTS idx_room_members_user_id;

DROP TABLE IF EXISTS notifications;

ALTER TABLE messages
    DROP COLUMN IF EXISTS metadata,
    DROP COLUMN IF EXISTS parent_message_id,
    DROP COLUMN IF EXISTS channel_id;

DROP TABLE IF EXISTS room_channels;
DROP TABLE IF EXISTS room_categories;
DROP TABLE IF EXISTS room_members;
-- +goose StatementEnd
