-- name: GenerateSession :one
INSERT INTO sessions (user_id_session, valid, user_agent) 
VALUES ($1, $2, $3) 
ON CONFLICT (user_id_session) DO
UPDATE SET updated_at = current_timestamp
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions WHERE user_id_session = $1;
