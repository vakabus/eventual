// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: event.sql

package gen

import (
	"context"
	"database/sql"
	"strings"
)

const addEventOrganizer = `-- name: AddEventOrganizer :exec
INSERT INTO event_organizers ( event_id, user_id ) VALUES ( ?1, ?2 )
`

type AddEventOrganizerParams struct {
	EventID int64
	UserID  int64
}

func (q *Queries) AddEventOrganizer(ctx context.Context, arg AddEventOrganizerParams) error {
	_, err := q.db.ExecContext(ctx, addEventOrganizer, arg.EventID, arg.UserID)
	return err
}

const addEventOrganizerByInviteCode = `-- name: AddEventOrganizerByInviteCode :one
INSERT INTO event_organizers ( event_id, user_id ) VALUES ((SELECT id FROM events WHERE invite_code = ?1), ?2) RETURNING id
`

type AddEventOrganizerByInviteCodeParams struct {
	InviteCode string
	UserID     int64
}

func (q *Queries) AddEventOrganizerByInviteCode(ctx context.Context, arg AddEventOrganizerByInviteCodeParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, addEventOrganizerByInviteCode, arg.InviteCode, arg.UserID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const addParticipant = `-- name: AddParticipant :one
INSERT INTO participants ( event_id, email, name) VALUES ( ?1, ?2, ?3 ) RETURNING id, event_id, email, name
`

type AddParticipantParams struct {
	EventID int64
	Email   string
	Name    sql.NullString
}

func (q *Queries) AddParticipant(ctx context.Context, arg AddParticipantParams) (Participant, error) {
	row := q.db.QueryRowContext(ctx, addParticipant, arg.EventID, arg.Email, arg.Name)
	var i Participant
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.Email,
		&i.Name,
	)
	return i, err
}

const addTemplate = `-- name: AddTemplate :one
INSERT INTO email_templates ( event_id, name, body ) VALUES ( ?1, ?2, ?3 ) RETURNING id, event_id, name, body
`

type AddTemplateParams struct {
	EventID int64
	Name    string
	Body    string
}

func (q *Queries) AddTemplate(ctx context.Context, arg AddTemplateParams) (EmailTemplate, error) {
	row := q.db.QueryRowContext(ctx, addTemplate, arg.EventID, arg.Name, arg.Body)
	var i EmailTemplate
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.Name,
		&i.Body,
	)
	return i, err
}

const eventOrganizers = `-- name: EventOrganizers :many
SELECT id, text_id, email, name, picture_url, deleted FROM users WHERE id IN (SELECT user_id FROM event_organizers WHERE event_id = ?1)
`

func (q *Queries) EventOrganizers(ctx context.Context, eventID int64) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, eventOrganizers, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.TextID,
			&i.Email,
			&i.Name,
			&i.PictureUrl,
			&i.Deleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const eventsByIds = `-- name: EventsByIds :many
SELECT e.id, e.name, e.description, e.invite_code, e.deleted FROM events e WHERE deleted = FALSE AND e.id IN (/*SLICE:event_ids*/?) AND e.id IN (SELECT eo.event_id FROM event_organizers eo WHERE eo.user_id = ?2)
`

type EventsByIdsParams struct {
	EventIds []int64
	UserID   int64
}

func (q *Queries) EventsByIds(ctx context.Context, arg EventsByIdsParams) ([]Event, error) {
	query := eventsByIds
	var queryParams []interface{}
	if len(arg.EventIds) > 0 {
		for _, v := range arg.EventIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:event_ids*/?", strings.Repeat(",?", len(arg.EventIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:event_ids*/?", "NULL", 1)
	}
	queryParams = append(queryParams, arg.UserID)
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.InviteCode,
			&i.Deleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listEvents = `-- name: ListEvents :many
SELECT id, name, description, invite_code, deleted FROM events WHERE deleted = FALSE AND id IN (SELECT event_id FROM event_organizers WHERE user_id = ?1)
`

func (q *Queries) ListEvents(ctx context.Context, userID int64) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, listEvents, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.InviteCode,
			&i.Deleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newEvent = `-- name: NewEvent :one
INSERT INTO events ( name, description, invite_code ) VALUES ( ?1, ?2, ?3 ) RETURNING id
`

type NewEventParams struct {
	Name        string
	Description string
	InviteCode  string
}

func (q *Queries) NewEvent(ctx context.Context, arg NewEventParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, newEvent, arg.Name, arg.Description, arg.InviteCode)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const participants = `-- name: Participants :many
SELECT id, event_id, email, name FROM participants WHERE event_id = ?1
`

func (q *Queries) Participants(ctx context.Context, eventID int64) ([]Participant, error) {
	rows, err := q.db.QueryContext(ctx, participants, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Participant
	for rows.Next() {
		var i Participant
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.Email,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeEventOrganizer = `-- name: RemoveEventOrganizer :exec
DELETE FROM event_organizers WHERE event_id = ?1 AND user_id = ?2
`

type RemoveEventOrganizerParams struct {
	EventID int64
	UserID  int64
}

func (q *Queries) RemoveEventOrganizer(ctx context.Context, arg RemoveEventOrganizerParams) error {
	_, err := q.db.ExecContext(ctx, removeEventOrganizer, arg.EventID, arg.UserID)
	return err
}

const removeParticipant = `-- name: RemoveParticipant :exec
DELETE FROM participants WHERE id = ?1 and event_id = ?2
`

type RemoveParticipantParams struct {
	ID      int64
	EventID int64
}

func (q *Queries) RemoveParticipant(ctx context.Context, arg RemoveParticipantParams) error {
	_, err := q.db.ExecContext(ctx, removeParticipant, arg.ID, arg.EventID)
	return err
}

const removeTemplate = `-- name: RemoveTemplate :exec
DELETE FROM email_templates WHERE id = ?1 AND event_id = ?2
`

type RemoveTemplateParams struct {
	ID      int64
	EventID int64
}

func (q *Queries) RemoveTemplate(ctx context.Context, arg RemoveTemplateParams) error {
	_, err := q.db.ExecContext(ctx, removeTemplate, arg.ID, arg.EventID)
	return err
}

const setInviteCode = `-- name: SetInviteCode :exec
UPDATE events SET invite_code = ?1 WHERE id = ?2
`

type SetInviteCodeParams struct {
	InviteCode string
	ID         int64
}

func (q *Queries) SetInviteCode(ctx context.Context, arg SetInviteCodeParams) error {
	_, err := q.db.ExecContext(ctx, setInviteCode, arg.InviteCode, arg.ID)
	return err
}

const template = `-- name: Template :one
SELECT id, event_id, name, body FROM email_templates WHERE id = ?1 AND event_id = ?2
`

type TemplateParams struct {
	ID      int64
	EventID int64
}

func (q *Queries) Template(ctx context.Context, arg TemplateParams) (EmailTemplate, error) {
	row := q.db.QueryRowContext(ctx, template, arg.ID, arg.EventID)
	var i EmailTemplate
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.Name,
		&i.Body,
	)
	return i, err
}

const templates = `-- name: Templates :many
SELECT id, event_id, name, body FROM email_templates WHERE event_id = ?1
`

func (q *Queries) Templates(ctx context.Context, eventID int64) ([]EmailTemplate, error) {
	rows, err := q.db.QueryContext(ctx, templates, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EmailTemplate
	for rows.Next() {
		var i EmailTemplate
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.Name,
			&i.Body,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateEvent = `-- name: UpdateEvent :one
UPDATE events SET name = ?1, description = ?2 WHERE id = ?3 RETURNING id, name, description, invite_code, deleted
`

type UpdateEventParams struct {
	Name        string
	Description string
	ID          int64
}

func (q *Queries) UpdateEvent(ctx context.Context, arg UpdateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, updateEvent, arg.Name, arg.Description, arg.ID)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.InviteCode,
		&i.Deleted,
	)
	return i, err
}

const updateParticipant = `-- name: UpdateParticipant :one
UPDATE participants SET email = ?1, name = ?2 WHERE id = ?3 and event_id = ?4 RETURNING id, event_id, email, name
`

type UpdateParticipantParams struct {
	Email   string
	Name    sql.NullString
	ID      int64
	EventID int64
}

func (q *Queries) UpdateParticipant(ctx context.Context, arg UpdateParticipantParams) (Participant, error) {
	row := q.db.QueryRowContext(ctx, updateParticipant,
		arg.Email,
		arg.Name,
		arg.ID,
		arg.EventID,
	)
	var i Participant
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.Email,
		&i.Name,
	)
	return i, err
}

const updateTemplate = `-- name: UpdateTemplate :exec
UPDATE email_templates SET name = ?1, body = ?2 WHERE id = ?3 AND event_id = ?4
`

type UpdateTemplateParams struct {
	Name    string
	Body    string
	ID      int64
	EventID int64
}

func (q *Queries) UpdateTemplate(ctx context.Context, arg UpdateTemplateParams) error {
	_, err := q.db.ExecContext(ctx, updateTemplate,
		arg.Name,
		arg.Body,
		arg.ID,
		arg.EventID,
	)
	return err
}
