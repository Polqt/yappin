-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_achievements (
    user_id UUID NOT NULL REFERENCES users(id),
    achievement_type_id UUID NOT NULL REFERENCES achievement_types(id),
    earned_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, achievement_type_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_achievements;
-- +goose StatementEnd