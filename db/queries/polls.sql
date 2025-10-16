-- name: CreatePoll :one
INSERT INTO polls (title, user_id)
VALUES ($1, $2)
RETURNING id,title, user_id;

-- name: GetPollByID :one
SELECT id, title, user_id
FROM polls
WHERE id = $1;

-- name: GetAllPolls :many
SELECT title, user_id
FROM polls;

-- name: UpdatePoll :exec
UPDATE polls
SET title = $2
WHERE id = $1;

-- name: DeletePoll :exec
DELETE FROM polls
WHERE id = $1;