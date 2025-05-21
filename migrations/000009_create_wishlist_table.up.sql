CREATE TABLE IF NOT EXISTS wishlist
(
    id           bigserial PRIMARY KEY,
    identity_id  uuid                        NOT NULL,
    furniture_id bigint REFERENCES furniture (id) ON DELETE CASCADE,
    created_at   timestamp(0) WITH TIME ZONE NOT NULL DEFAULT now(),
    UNIQUE (identity_id, furniture_id)
);