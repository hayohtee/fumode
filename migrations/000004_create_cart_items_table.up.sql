CREATE TABLE IF NOT EXISTS cart_items
(
    id           bigserial PRIMARY KEY,
    cart_id      bigint    NOT NULL REFERENCES cart (id) ON DELETE cascade,
    furniture_id bigserial NOT NULL REFERENCES furniture (id),
    quantity     int       NOT NULL
);