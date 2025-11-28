-- name: CreatePoll :one
INSERT INTO polls (title, user_id)
VALUES (@title, @user_id)
RETURNING id,title, user_id;

-- name: GetPollByID :many
SELECT 
    polls.id,
    polls.title,
    polls.user_id,
    options.id AS option_id,
    options.content AS option_content
FROM polls
inner JOIN options ON polls.id = options.poll_id
WHERE polls.id = @id
ORDER BY polls.id ASC;

-- name: GetAllPolls :many
SELECT
    p.id AS poll_id,
    p.title,
    p.user_id,
    o.id AS option_id,
    o.content AS option_content,
    r.option_id AS user_voted_option_id
FROM polls p
JOIN options o ON p.id = o.poll_id
LEFT JOIN results r ON p.id = r.poll_id AND r.user_id = @user_id
ORDER BY p.id ASC;

-- name: UpdatePoll :exec
UPDATE polls
SET title = @title
WHERE id = @id;

-- name: DeletePoll :exec
DELETE FROM polls
WHERE id = @id;

-- name: GetPollsByUserID :many
SELECT
    p.id AS poll_id,
    p.title,
    p.user_id,
    o.id AS option_id,
    o.content AS option_content,
    r.option_id AS user_voted_option_id
FROM polls p
JOIN options o ON p.id = o.poll_id
LEFT JOIN results r ON p.id = r.poll_id AND r.user_id = @viewer_id
WHERE p.user_id = @owner_id
ORDER BY p.id ASC;