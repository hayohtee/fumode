CREATE TABLE IF NOT EXISTS reviews
(
    id           bigserial PRIMARY KEY,
    identity_id  uuid   NOT NULL UNIQUE,
    furniture_id bigint NOT NULL REFERENCES furniture (id) ON DELETE cascade,
    rating       int    NOT NULL,
    comment      text   NOT NULL,
    created_at   timestamp(0) with time zone DEFAULT now(),
    version      int    NOT NULL             DEFAULT 1,

    CONSTRAINT check_rating_value CHECK ( rating >= 1 AND rating <= 5),
    UNIQUE (identity_id, furniture_id)
);