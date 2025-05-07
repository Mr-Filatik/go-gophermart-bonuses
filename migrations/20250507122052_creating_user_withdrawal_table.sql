-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_withdrawals (
    "id" SERIAL PRIMARY KEY,
    "order" TEXT UNIQUE NOT NULL,
    "sum" DOUBLE PRECISION NOT NULL,
    "processed_at" TIMESTAMP NOT NULL,
    "user_id" INTEGER NOT NULL,

    CONSTRAINT fk_user
        FOREIGN KEY ("user_id")
        REFERENCES users("id")
        ON DELETE CASCADE
);
COMMENT ON TABLE user_withdrawals IS 'User withdrawals table';
COMMENT ON COLUMN user_withdrawals.id IS 'Unique withdrawal identifier';
COMMENT ON COLUMN user_withdrawals.order IS 'Unique order number';
COMMENT ON COLUMN user_withdrawals.sum IS 'Withdrawal amount';
COMMENT ON COLUMN user_withdrawals.processed_at IS 'Order processed time';
COMMENT ON COLUMN user_withdrawals.user_id IS 'User identifier (Foreign Key)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_withdrawals;
-- +goose StatementEnd
