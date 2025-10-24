-- name: CreateResult :exec
INSERT INTO results (poll_id, option_id, user_id)
VALUES (@poll_id, @option_id, @user_id);

-- name: GetResultByID :one
SELECT id, poll_id, option_id, user_id
FROM results
WHERE id = @id;

-- name: GetAllResults :many
SELECT id, poll_id, option_id, user_id
FROM results
ORDER BY id ASC;

-- name: UpdateResult :exec
UPDATE results
SET option_id = @option_id
WHERE id = @id;

-- name: DeleteResult :exec
DELETE FROM results
WHERE id = @id;

-- name: GetResultsByPollID :many
SELECT id, poll_id, option_id, user_id
FROM results
WHERE poll_id = @poll_id;

-- name: GetResultsGroupByPollID :many
SELECT poll_id, option_id, COUNT(*) AS total
FROM results
GROUP BY poll_id, option_id;
