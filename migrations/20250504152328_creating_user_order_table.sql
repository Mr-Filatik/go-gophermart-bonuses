-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TABLE IF NOT EXISTS user_orders (
    id SERIAL PRIMARY KEY,
    number INTEGER NOT NULL UNIQUE,
    uploaded_at TIMESTAMP NOT NULL,
    status user_order_status NOT NULL,
    user_id INTEGER NOT NULL,

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
COMMENT ON TABLE user_orders IS 'User orders table';
COMMENT ON COLUMN user_orders.id IS 'Unique order identifier';
COMMENT ON COLUMN user_orders.number IS 'Unique order number';
COMMENT ON COLUMN user_orders.uploaded_at IS 'Order upload time';
COMMENT ON COLUMN user_orders.status IS 'Order status (NEW, PROCESSING, INVALID, PROCESSED)';
COMMENT ON COLUMN user_orders.user_id IS 'User identifier (Foreign Key)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_orders;
DROP TYPE user_order_status;
-- +goose StatementEnd
