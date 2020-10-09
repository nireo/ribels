-- Your SQL goes here
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    osu_id TEXT UNIQUE NOT NULL,
    discord_id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP
);
