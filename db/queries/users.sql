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
    username = $2,
    email    = $3,
    password = $4
WHERE id = $1
RETURNING id, username, email;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING username;