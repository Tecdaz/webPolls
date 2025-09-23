-- name: CreateUser :one
INSERT INTO users (username,password, email)
VALUES ($1, $2, $3)
RETURNING id, username, email;

-- name: GetUserByID :one
SELECT id, username, email
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT id, username, email
FROM users
WHERE username = $1;

-- name: GetAllUsers :many
SELECT id, username, email
FROM users;

-- name: UpdateUser :exec
UPDATE users
SET username = $2, password = $3, email = $4
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;