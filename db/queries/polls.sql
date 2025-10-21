-- name: CreatePoll :one
INSERT INTO polls (title, user_id)
VALUES (@title, @user_id)
RETURNING id,title, user_id;

-- name: GetPollByID :one
SELECT id, title, user_id
FROM polls
WHERE id = @id;

-- name: GetAllPolls :many
SELECT
    p.id AS poll_id,
    p.title,
    p.user_id,
    o.id AS option_id,
    o.content,
    o.correct
FROM polls p
LEFT JOIN options o ON p.id = o.poll_id;


-- name: UpdatePoll :exec
UPDATE polls
SET title = @title
WHERE id = @id;

-- name: DeletePoll :exec
DELETE FROM polls
WHERE id = @id;