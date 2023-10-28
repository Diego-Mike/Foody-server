// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (social_id, username, email, picture, provider) VALUES ($1, $2, $3, $4, $5) RETURNING user_id, social_id, username, email, picture, provider, registered_at
`

type CreateUserParams struct {
	SocialID string `json:"social_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.SocialID,
		arg.Username,
		arg.Email,
		arg.Picture,
		arg.Provider,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.SocialID,
		&i.Username,
		&i.Email,
		&i.Picture,
		&i.Provider,
		&i.RegisteredAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT user_id, social_id, username, email, picture, provider, registered_at FROM users WHERE user_id = $1
`

func (q *Queries) GetUserById(ctx context.Context, userID int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.SocialID,
		&i.Username,
		&i.Email,
		&i.Picture,
		&i.Provider,
		&i.RegisteredAt,
	)
	return i, err
}

const getUserBySocialId = `-- name: GetUserBySocialId :one
SELECT user_id, social_id, username, email, picture, provider, registered_at FROM users WHERE social_id = $1
`

func (q *Queries) GetUserBySocialId(ctx context.Context, socialID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserBySocialId, socialID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.SocialID,
		&i.Username,
		&i.Email,
		&i.Picture,
		&i.Provider,
		&i.RegisteredAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET username = $1::varchar, email = $2::varchar, picture = $3::varchar
WHERE users.user_id = $4::bigint 
AND (username <> $5::varchar OR email <> $6::varchar OR picture <> $7::varchar) 
RETURNING user_id, social_id, username, email, picture, provider, registered_at
`

type UpdateUserParams struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Pictue      string `json:"pictue"`
	UserID      int64  `json:"user_id"`
	OldUsername string `json:"old_username"`
	OldEmail    string `json:"old_email"`
	OldPicture  string `json:"old_picture"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Username,
		arg.Email,
		arg.Pictue,
		arg.UserID,
		arg.OldUsername,
		arg.OldEmail,
		arg.OldPicture,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.SocialID,
		&i.Username,
		&i.Email,
		&i.Picture,
		&i.Provider,
		&i.RegisteredAt,
	)
	return i, err
}