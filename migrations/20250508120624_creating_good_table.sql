-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goods (
    id BIGSERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    price BIGINT NOT NULL,
    order_id BIGINT NOT NULL,

    CONSTRAINT fk_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE
);
COMMENT ON TABLE goods IS 'Goods table';
COMMENT ON COLUMN goods.id IS 'Unique good identifier';
COMMENT ON COLUMN goods.description IS 'Good description';
COMMENT ON COLUMN goods.price IS 'Good price';
COMMENT ON COLUMN goods.order_id IS 'Order identifier (Foreign Key)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goods;
-- +goose StatementEnd
