CREATE TABLE characters (
    id   BIGSERIAL PRIMARY KEY,
    name text      NOT NULL,
    bio  text      NOT NULL
);

ALTER TABLE characters ADD COLUMN note text;
