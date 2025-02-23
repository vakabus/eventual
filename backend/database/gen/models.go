// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package gen

import (
	"database/sql"
)

type EmailTemplate struct {
	ID      int64
	EventID int64
	Name    string
	Body    string
}

type Event struct {
	ID          int64
	Name        string
	Description string
	InviteCode  string
	Deleted     sql.NullInt64
}

type EventOrganizer struct {
	ID      int64
	EventID int64
	UserID  int64
}

type Participant struct {
	ID      int64
	EventID int64
	Json    string
}

type Session struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt string
}

type User struct {
	ID         int64
	TextID     string
	Email      string
	Name       string
	PictureUrl string
	Deleted    sql.NullInt64
}
