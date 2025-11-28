-- name: Vote :exec
INSERT INTO results (poll_id, option_id, user_id)
VALUES (@poll_id, @option_id, @user_id)
ON CONFLICT (poll_id, option_id, user_id) DO NOTHING;

-- name: DeleteUserVote :exec
DELETE FROM results
WHERE poll_id = @poll_id AND user_id = @user_id;

-- name: VoteOneStep :exec
WITH deleted AS (
    DELETE FROM results 
    WHERE poll_id = @poll_id AND user_id = @user_id
)
INSERT INTO results (poll_id, option_id, user_id)
VALUES (@poll_id, @option_id, @user_id);

-- name: GetPollResults :many
SELECT 
    option_id,
    COUNT(user_id) AS vote_count
FROM results
WHERE poll_id = @poll_id
GROUP BY option_id;

-- name: GetUserVote :one
SELECT option_id
FROM results
WHERE poll_id = @poll_id AND user_id = @user_id;