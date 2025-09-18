-- name: CreateOption :one
INSERT INTO options (content, correct)
VALUES ($1, $2)
RETURNING content, correct ;

-- name: GetOptionByID :one
SELECT content, correct
FROM options
WHERE id = $1;

-- name: GetAllOptions :many
SELECT content, correct
FROM options;

-- name: UpdateOption :exec
UPDATE options
SET content = $2, correct = $3
WHERE id = $1;

-- name: DeleteOption :exec
DELETE FROM options
WHERE id = $1;