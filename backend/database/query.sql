-- name: UserByName :one
SELECT * FROM users WHERE name = :name LIMIT 1;