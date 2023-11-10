-- name: GetUserBySocialId :one
SELECT * FROM users WHERE social_id = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE user_id = $1;

-- name: GetFullUser :one
SELECT u.user_id, u.social_id, u.username, u.email, u.picture, bm.business_id, bm.business_position FROM users u 
LEFT JOIN business_members bm ON u.user_id  = bm.user_id WHERE u.user_id = $1;

-- name: CreateUser :one
INSERT INTO users (social_id, username, email, picture, provider) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET username = sqlc.arg(username)::varchar, email = sqlc.arg(email)::varchar, picture = sqlc.arg(pictue)::varchar
WHERE users.user_id = sqlc.arg(user_id)::bigint 
AND (username <> sqlc.arg(old_username)::varchar OR email <> sqlc.arg(old_email)::varchar OR picture <> sqlc.arg(old_picture)::varchar) 
RETURNING *;
