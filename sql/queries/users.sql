-- name: CreateUser :one
INSERT INTO
    users (user_name, email, password, created_at, updated_at)
VALUES
($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUsers :many
SELECT
    id,
    user_name,
    email,
    created_at,
    updated_at,
    deleted_at
FROM users
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: GetUserByID :one
SELECT
    id,
    user_name,
    email,
    created_at,
    updated_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByName :one
SELECT
    id,
    user_name,
    email,
    created_at,
    updated_at
FROM users
WHERE user_name = $1 AND deleted_at IS NULL;

-- name: GetUserByEmailForAuth :one
SELECT id, email, password
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE users
SET
    user_name = $2,
    email = $3,
    updated_at = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET
    deleted_at = $2
WHERE id = $1;

-- name: RestoreUser :exec
UPDATE users
SET
    deleted_at = NULL
WHERE id = $1;

-- name: DestroyUser :exec
DELETE FROM users
WHERE id = $1;
