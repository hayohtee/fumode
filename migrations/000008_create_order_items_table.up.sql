CREATE TABLE IF NOT EXISTS order_items
(
    id           bigserial PRIMARY KEY,
    order_id     bigint REFERENCES orders (id) ON DELETE CASCADE,
    furniture_id bigserial REFERENCES furniture (id) ON DELETE CASCADE,
    quantity     int                         NOT NULL CHECK ( quantity > 0 ),
    price        decimal(10, 2)              NOT NULL,
    created_at   timestamp(0) with time zone NOT NULL DEFAULT now(),
    UNIQUE (order_id, furniture_id)
);