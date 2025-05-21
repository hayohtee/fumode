CREATE TABLE IF NOT EXISTS furniture
(
    id             bigserial PRIMARY KEY,
    name           text                        NOT NULL,
    description    text                        NOT NULL,
    price          decimal(10, 2)              NOT NULL,
    stock_quantity int                         NOT NULL,
    category_id    bigint REFERENCES categories (id) ON UPDATE cascade,
    created_at     timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at     timestamp(0) with time zone NOT NULL DEFAULT now(),
    version        int                         NOT NULL DEFAULT 1
);