-- name: ListEvents :many
SELECT * FROM events WHERE deleted = FALSE AND id IN (SELECT event_id FROM event_organizers WHERE user_id = :user_id);

-- name: EventsByIds :many
SELECT e.* FROM events e WHERE deleted = FALSE AND e.id IN (sqlc.slice(event_ids)) AND e.id IN (SELECT eo.event_id FROM event_organizers eo WHERE eo.user_id = :user_id);

-- name: NewEvent :one
INSERT INTO events ( name, description, invite_code ) VALUES ( :name, :description, :invite_code ) RETURNING id;

-- name: SetInviteCode :exec
UPDATE events SET invite_code = :invite_code WHERE id = :id;

-- name: UpdateEvent :one
UPDATE events SET name = :name, description = :description WHERE id = :id RETURNING *;



-- name: AddEventOrganizer :exec
INSERT INTO event_organizers ( event_id, user_id ) VALUES ( :event_id, :user_id );

-- name: AddEventOrganizerByInviteCode :one
INSERT INTO event_organizers ( event_id, user_id ) VALUES ((SELECT id FROM events WHERE invite_code = :invite_code), :user_id) RETURNING id;

-- name: EventOrganizers :many
SELECT * FROM users WHERE id IN (SELECT user_id FROM event_organizers WHERE event_id = :event_id);

-- name: RemoveEventOrganizer :exec
DELETE FROM event_organizers WHERE event_id = :event_id AND user_id = :user_id;



-- name: AddParticipant :one
INSERT INTO participants ( event_id, email ) VALUES ( :event_id, :email ) RETURNING *;

-- name: Participants :many
SELECT * FROM participants WHERE event_id = :event_id;

-- name: RemoveParticipant :exec
DELETE FROM participants WHERE id = :id;