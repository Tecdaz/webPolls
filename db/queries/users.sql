-- name: CreateUser :one
INSERT INTO users (username,password, email)
VALUES (@username, @password, @email)
RETURNING id, username, email;

-- name: GetUserByID :one
SELECT id, username, email
FROM users
WHERE id = @id;

-- name: GetUserByUsername :one
SELECT id, username, email
FROM users
WHERE username = @username;

-- name: GetUserByEmail :one
SELECT id, username, email
FROM users
WHERE email = @email;

-- name: GetAllUsers :many
SELECT id, username, email
FROM users;

-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE(sqlc.narg(username), username),
    email    = COALESCE(sqlc.narg(email), email),
    password = COALESCE(sqlc.narg(password), password)
WHERE id = @id
RETURNING id, username, email;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = @id
RETURNING username;