-- name: CreatePoll :one
INSERT INTO polls (title, user_id)
VALUES (@title, @user_id)
RETURNING id,title, user_id;

-- name: GetPollByID :one
SELECT id, title, user_id
FROM polls
WHERE id = @id;

-- name: GetAllPolls :many
SELECT title, user_id
FROM polls;

-- name: UpdatePoll :exec
UPDATE polls
SET title = @title
WHERE id = @id;

-- name: DeletePoll :exec
DELETE FROM polls
WHERE id = @id;