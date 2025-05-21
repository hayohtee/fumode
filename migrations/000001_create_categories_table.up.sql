CREATE TABLE IF NOT EXISTS categories
(
    id         bigserial PRIMARY KEY,
    name       text                        NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);