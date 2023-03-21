-- name: GetCharacter :one
SELECT * FROM characters
WHERE id = $1 LIMIT 1;

-- name: FindCharacter :one
SELECT * FROM characters
WHERE name = $1 LIMIT 1;

-- name: ListCharacter :many
SELECT * FROM characters
ORDER BY name;

-- name: CreateCharacter :one
INSERT INTO characters (
    name, bio, note
) VALUES (
             $1, $2, $3
         )
    RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE id = $1;