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