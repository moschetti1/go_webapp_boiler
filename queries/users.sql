-- name: GetUserByGoogleId :one
SELECT * FROM users
WHERE google_id = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, google_id, email, name, avatar_url)
VALUES (?, ?, ?, ?, ?)
RETURNING *;
