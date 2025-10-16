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

-- name: GetUserByEmail :one
SELECT id, username, email
FROM users
WHERE email = $1;

-- name: GetAllUsers :many
SELECT id, username, email
FROM users;

-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE(NULLIF($2, ''), username),
    email    = COALESCE(NULLIF($3, ''), email),
    password = COALESCE(NULLIF($4, ''), password)
WHERE id = $1
RETURNING id, username, email;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;