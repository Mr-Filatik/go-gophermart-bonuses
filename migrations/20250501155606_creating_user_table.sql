-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);
COMMENT ON TABLE users IS 'Users table';
COMMENT ON COLUMN users.id IS 'Unique user identifier';
COMMENT ON COLUMN users.login IS 'Unique user login';
COMMENT ON COLUMN users.password_hash IS 'User password hash';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
