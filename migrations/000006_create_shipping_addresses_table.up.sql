CREATE TABLE IF NOT EXISTS shipping_addresses
(
    id          bigserial PRIMARY KEY,
    identity_id uuid                        NOT NULL UNIQUE,
    address     text                        NOT NULL,
    city        text                        NOT NULL,
    state       text                        NOT NULL,
    country     text                        NOT NULL,
    postal_code int                         NOT NULL,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT now()
);