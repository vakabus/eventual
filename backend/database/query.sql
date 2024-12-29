-- name: UserByTextId :one
SELECT * FROM users WHERE text_id = :text_id LIMIT 1;

-- name: NewUser :one
INSERT INTO users ( text_id, email, name, picture_url ) VALUES ( :text_id, :email, :name, :picture_url ) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = :email, name = :name, picture_url = :picture_url WHERE text_id = :text_id RETURNING *;

-- name: AddSession :exec
INSERT INTO sessions ( user_id, token, expires_at ) VALUES ( :user_id, :token, :expires_at );

-- name: CleanSessions :exec
DELETE FROM sessions WHERE expires_at < datetime('now');

-- name: UserBySession :one
SELECT * FROM users WHERE id = (SELECT user_id FROM sessions WHERE token = :token AND expires_at > datetime('now')) LIMIT 1;



-- name: ListEvents :many
SELECT * FROM events WHERE deleted = FALSE AND id IN (SELECT event_id FROM event_organizers WHERE user_id = :user_id);

-- name: NewEvent :one
INSERT INTO events ( name, description ) VALUES ( :name, :description ) RETURNING id;

-- name: AddEventOrganizer :exec
INSERT INTO event_organizers ( event_id, user_id ) VALUES ( :event_id, :user_id );