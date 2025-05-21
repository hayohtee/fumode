CREATE TABLE IF NOT EXISTS categories
(
    id         bigserial PRIMARY KEY,
    name       TEXT                        NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);