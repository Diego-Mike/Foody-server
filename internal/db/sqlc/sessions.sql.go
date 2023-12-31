// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: sessions.sql

package db

import (
	"context"
)

const generateSession = `-- name: GenerateSession :one
INSERT INTO sessions (user_id_session, valid, user_agent) 
VALUES ($1, $2, $3) 
ON CONFLICT (user_id_session) DO
UPDATE SET updated_at = current_timestamp
RETURNING user_id_session, valid, user_agent, created_at, updated_at
`

type GenerateSessionParams struct {
	UserIDSession int64  `json:"user_id_session"`
	Valid         bool   `json:"valid"`
	UserAgent     string `json:"user_agent"`
}

func (q *Queries) GenerateSession(ctx context.Context, arg GenerateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, generateSession, arg.UserIDSession, arg.Valid, arg.UserAgent)
	var i Session
	err := row.Scan(
		&i.UserIDSession,
		&i.Valid,
		&i.UserAgent,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT user_id_session, valid, user_agent, created_at, updated_at FROM sessions WHERE user_id_session = $1
`

func (q *Queries) GetSession(ctx context.Context, userIDSession int64) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, userIDSession)
	var i Session
	err := row.Scan(
		&i.UserIDSession,
		&i.Valid,
		&i.UserAgent,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
