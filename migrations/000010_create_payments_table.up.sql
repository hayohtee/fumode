CREATE TABLE IF NOT EXISTS payments
(
    id                 bigserial PRIMARY KEY,
    identity_id        uuid                        NOT NULL,
    order_id           bigint                      REFERENCES orders (id) ON DELETE SET NULL,
    amount             decimal(10, 2)              NOT NULL,
    provider           text                        NOT NULL,
    provider_reference text                        NOT NULL,
    paid_at            timestamp(0) with time zone,
    status             text                        NOT NULL CHECK ( status IN ('success', 'failed', 'pending') ),
    created_at         timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at         timestamp(0) with time zone NOT NULL DEFAULT now()
);