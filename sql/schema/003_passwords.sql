-- +goose Up
ALTER TABLE users
ADD COLUMN password TEXT DEFAULT 'unset' NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN password;