CREATE TABLE IF NOT EXISTS cart
(
    id          bigserial PRIMARY KEY,
    identity_id uuid                        NOT NULL,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT now()
);