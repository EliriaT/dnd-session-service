CREATE TABLE "sessions" (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    campaign_id BIGINT NOT NULL,
    map_id BIGINT NOT NULL
);