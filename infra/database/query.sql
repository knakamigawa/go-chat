-- name: GetCharacter :one
SELECT *
FROM characters
WHERE id = $1 LIMIT 1;

-- name: FindCharacter :one
SELECT *
FROM characters
WHERE name = $1 LIMIT 1;

-- name: ListCharacter :many
SELECT *
FROM characters
ORDER BY name;

-- name: CreateCharacter :one
INSERT INTO characters (name, bio, note)
VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteCharacter :exec
DELETE
FROM characters
WHERE id = $1;

-- name: FindUserByEmail :one
SELECT users.id, users.login_name, user_login_with_email_passwords.email, user_login_with_email_passwords.password_hash
FROM users
         JOIN user_login_with_email_passwords
              ON users.id = user_login_with_email_passwords.user_id
WHERE user_login_with_email_passwords.email = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, login_name)
VALUES ($1, $2) RETURNING *;

-- name: CreateUserEmailPassword :one
INSERT INTO user_login_with_email_passwords (user_id, email, password_hash)
VALUES ($1, $2, $3) RETURNING *;