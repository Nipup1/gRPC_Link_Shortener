-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links(
    id SERIAL PRIMARY KEY,
    link TEXT NOT NULL UNIQUE,
    short_link TEXT NOT NULL UNIQUE);
CREATE INDEX IF NOT EXISTS idx_short_link ON links(short_link);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_short_link;
DROP TABLE IF EXISTS links;
-- +goose StatementEnd
