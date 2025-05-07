-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN current DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN withdrawn DOUBLE PRECISION NOT NULL DEFAULT 0;
COMMENT ON COLUMN users.current IS 'Current balance of the user';
COMMENT ON COLUMN users.withdrawn IS 'Total amount withdrawn by the user';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN current;
ALTER TABLE users DROP COLUMN withdrawn;
-- +goose StatementEnd
