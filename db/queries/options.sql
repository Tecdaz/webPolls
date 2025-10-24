-- name: CreateOption :one
INSERT INTO options (content, poll_id)
VALUES (@content, @poll_id)
RETURNING id, content, poll_id;
 

-- name: GetAllOptions :many
SELECT id, content, poll_id
FROM options
ORDER BY id ASC;

-- name: UpdateOption :one
UPDATE options
SET content = @content
WHERE id = @id
RETURNING id, content, poll_id;

-- name: DeleteOption :exec
DELETE FROM options
WHERE id = @id;

-- name: GetOptionByID :one
SELECT id, content, poll_id
FROM options
WHERE id = @id;

-- name: GetOptionByPollID :many
SELECT id, content, poll_id
FROM options
WHERE poll_id = @poll_id
ORDER BY id ASC;