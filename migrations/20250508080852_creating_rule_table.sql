-- +goose Up
-- +goose StatementBegin
CREATE TYPE reward_types AS ENUM ('%', 'pt');
CREATE TABLE IF NOT EXISTS rules (
    id SERIAL PRIMARY KEY,
    match TEXT NOT NULL,
    reward DOUBLE PRECISION NOT NULL,
    reward_type reward_types NOT NULL
);
COMMENT ON TABLE rules IS 'Rules table';
COMMENT ON COLUMN rules.id IS 'Unique rule identifier';
COMMENT ON COLUMN rules.match IS 'Match name';
COMMENT ON COLUMN rules.reward IS 'Amount of reward';
COMMENT ON COLUMN rules.reward_type IS 'Reward type (%, pt)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rules;
DROP TYPE reward_types;
-- +goose StatementEnd
