-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS achievement_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    threshold_type VARCHAR(50) NOT NULL,
    threshold_value INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS achievement_types;
-- +goose StatementEnd