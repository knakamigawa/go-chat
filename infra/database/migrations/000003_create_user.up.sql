CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    login_name text NOT NULL UNIQUE
);

CREATE TABLE user_login_with_email_passwords
(
    id            BIGSERIAL PRIMARY KEY,
    user_id       UUID      NOT NULL REFERENCES users (id),
    email         text NOT NULL UNIQUE,
    password_hash text NOT NULL
);
