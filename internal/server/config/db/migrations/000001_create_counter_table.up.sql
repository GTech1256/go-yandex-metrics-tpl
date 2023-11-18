CREATE TABLE counter (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    delta BIGINT NOT NULL
);