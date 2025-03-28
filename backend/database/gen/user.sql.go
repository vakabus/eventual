// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package gen

import (
	"context"
)

const addSession = `-- name: AddSession :exec
INSERT INTO sessions ( user_id, token, expires_at ) VALUES ( ?1, ?2, ?3 )
`

type AddSessionParams struct {
	UserID    int64
	Token     string
	ExpiresAt string
}

func (q *Queries) AddSession(ctx context.Context, arg AddSessionParams) error {
	_, err := q.db.ExecContext(ctx, addSession, arg.UserID, arg.Token, arg.ExpiresAt)
	return err
}

const cleanSessions = `-- name: CleanSessions :exec
DELETE FROM sessions WHERE expires_at < datetime('now')
`

func (q *Queries) CleanSessions(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, cleanSessions)
	return err
}

const newUser = `-- name: NewUser :one
INSERT INTO users ( text_id, email, name, picture_url ) VALUES ( ?1, ?2, ?3, ?4 ) RETURNING id, text_id, email, name, picture_url, deleted
`

type NewUserParams struct {
	TextID     string
	Email      string
	Name       string
	PictureUrl string
}

func (q *Queries) NewUser(ctx context.Context, arg NewUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, newUser,
		arg.TextID,
		arg.Email,
		arg.Name,
		arg.PictureUrl,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TextID,
		&i.Email,
		&i.Name,
		&i.PictureUrl,
		&i.Deleted,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET email = ?1, name = ?2, picture_url = ?3 WHERE text_id = ?4 RETURNING id, text_id, email, name, picture_url, deleted
`

type UpdateUserParams struct {
	Email      string
	Name       string
	PictureUrl string
	TextID     string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Email,
		arg.Name,
		arg.PictureUrl,
		arg.TextID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TextID,
		&i.Email,
		&i.Name,
		&i.PictureUrl,
		&i.Deleted,
	)
	return i, err
}

const userBySession = `-- name: UserBySession :one
SELECT id, text_id, email, name, picture_url, deleted FROM users WHERE id = (SELECT user_id FROM sessions WHERE token = ?1 AND expires_at > datetime('now')) LIMIT 1
`

func (q *Queries) UserBySession(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, userBySession, token)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TextID,
		&i.Email,
		&i.Name,
		&i.PictureUrl,
		&i.Deleted,
	)
	return i, err
}

const userByTextId = `-- name: UserByTextId :one
SELECT id, text_id, email, name, picture_url, deleted FROM users WHERE text_id = ?1 LIMIT 1
`

func (q *Queries) UserByTextId(ctx context.Context, textID string) (User, error) {
	row := q.db.QueryRowContext(ctx, userByTextId, textID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TextID,
		&i.Email,
		&i.Name,
		&i.PictureUrl,
		&i.Deleted,
	)
	return i, err
}
