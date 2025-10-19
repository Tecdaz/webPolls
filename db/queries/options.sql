-- name: CreateOption :one
INSERT INTO options (content, correct, poll_id)
VALUES (@content, @correct, @poll_id)
RETURNING id, content, correct, poll_id;
 

-- name: GetOptionByID :one
SELECT content, correct
FROM options
WHERE id = @id;

-- name: GetOptionByPollID :many
SELECT content, correct
FROM options
WHERE poll_id = @poll_id;

-- name: GetAllOptions :many
SELECT content, correct
FROM options;

-- name: UpdateOption :exec
UPDATE options
SET content = @content, correct = @correct
WHERE id = @id;

-- name: DeleteOption :exec
DELETE FROM options
WHERE id = @id;