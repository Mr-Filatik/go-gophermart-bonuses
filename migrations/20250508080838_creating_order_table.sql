-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('REGISTERED', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    number INTEGER NOT NULL UNIQUE,
    uploaded_at TIMESTAMP NOT NULL,
    status order_status NOT NULL,
    accrual DOUBLE PRECISION NOT NULL
);
COMMENT ON TABLE orders IS 'Orders table';
COMMENT ON COLUMN orders.id IS 'Unique order identifier';
COMMENT ON COLUMN orders.number IS 'Unique order number';
COMMENT ON COLUMN orders.uploaded_at IS 'Order upload time';
COMMENT ON COLUMN orders.status IS 'Order status (REGISTERED, PROCESSING, INVALID, PROCESSED)';
COMMENT ON COLUMN orders.accrual IS 'Amount of reward';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE order_status;
-- +goose StatementEnd
