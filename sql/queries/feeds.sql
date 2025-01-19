-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name,url, user_id )
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

/*
-- name: GetOneUserByName :one
SELECT * FROM users 
WHERE users.name=$1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;
*/