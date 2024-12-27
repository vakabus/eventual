-- name: UserByTextId :one
SELECT * FROM users WHERE text_id = :text_id LIMIT 1;

-- name: NewUser :exec
INSERT INTO users ( text_id, email, name, picture_url ) VALUES ( :text_id, :email, :name, :picture_url );

-- name: UpdateUser :exec
UPDATE users SET email = :email, name = :name, picture_url = :picture_url WHERE text_id = :text_id;



-- name: ListEvents :many
SELECT * FROM events WHERE deleted = FALSE AND id IN (SELECT event_id FROM event_organizers WHERE user_id = :user_id);