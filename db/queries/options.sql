-- name: CreateOption :one
INSERT INTO options (content, correct, poll_id)
VALUES (@content, @correct, @poll_id)
RETURNING id, content, correct, poll_id;
 

-- name: GetAllOptions :many
SELECT id, content, correct
FROM options
ORDER BY id ASC;

-- name: UpdateOption :exec
UPDATE options
SET content = @content, correct = @correct
WHERE id = @id;

-- name: DeleteOption :exec
DELETE FROM options
WHERE id = @id;

-- name: GetOptionByID :one
SELECT id, content, correct, poll_id
FROM options
WHERE id = @id;

-- name: GetOptionByPollID :many
SELECT id, content, correct
FROM options
WHERE poll_id = @poll_id
ORDER BY id ASC;

-- name: UnsetOtherOptionsCorrect :exec
UPDATE options
SET correct = false
WHERE poll_id = @poll_id AND id != @id;