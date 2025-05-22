CREATE TABLE IF NOT EXISTS furniture_images
(
    id           bigserial PRIMARY KEY,
    furniture_id bigint                      NOT NULL REFERENCES furniture (id) ON DELETE CASCADE,
    url          text                        NOT NULL,
    is_primary   boolean                     NOT NULL DEFAULT false,
    created_at   timestamp(0) WITH TIME ZONE NOT NULL DEFAULT now()
);