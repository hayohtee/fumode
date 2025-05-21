CREATE TABLE IF NOT EXISTS orders
(
    id                  bigserial PRIMARY KEY,
    identity_id         uuid           NOT NULL,
    shipping_address_id bigint         NOT NULL REFERENCES shipping_addresses (id),
    status              text           NOT NULL CHECK ( status IN ('pending', 'shipped', 'delivered', 'canceled') ) DEFAULT 'pending',
    total_amount        decimal(10, 2) NOT NULL,
    payment_status      text CHECK ( payment_status IN ('unpaid', 'paid', 'failed') )                               DEFAULT 'unpaid',
    created_at          timestamp(0) with time zone                                                                 DEFAULT now(),
    updated_at          timestamp(0) with time zone                                                                 DEFAULT now(),
    version             int            NOT NULL                                                                     DEFAULT 1
);