-- name: CreateLink :one
INSERT INTO links (user_id, target_url, shortened_url, password, status, expiration_date, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetLinks :many
SELECT
    l.id,
    u.user_name,
    l.target_url,
    l.shortened_url,
    l.status,
    l.expiration_date,
    l.created_at,
    l.updated_at
FROM links l
JOIN users u ON l.user_id = u.id
WHERE deleted_at IS NULL;

-- name: GetLinkByID :one
SELECT
    l.id,
    u.user_name,
    l.target_url,
    l.shortened_url,
    l.status,
    l.expiration_date,
    l.created_at,
    l.updated_at
FROM links l
JOIN users u ON l.user_id = u.id
WHERE l.id = $1 AND delete_at IS NULL;

-- name: GetLinkByTargetURL :one
SELECT
    l.id,
    u.user_name,
    l.target_url,
    l.shortened_url,
    l.status,
    l.expiration_date,
    l.created_at,
    l.updated_at
FROM links l
JOIN users u ON l.user_id = u.id
WHERE target_url = $1 AND deleted_at IS NULL;

-- name: GetLinkByShortenedURL :one
SELECT
    l.id,
    u.user_name,
    l.target_url,
    l.shortened_url,
    l.status,
    l.expiration_date,
    l.created_at,
    l.updated_at
FROM links l
JOIN users u ON l.user_id = u.id
WHERE shortened_url = $1 AND deleted_at IS NULL;

-- name: GetTargetURLByShortenedURL :one
SELECT
    target_url
FROM links
WHERE shortened_url = $1 AND deleted_at IS NULL;

-- name: GetLinksByStatus :many
SELECT
    l.id,
    u.user_name,
    l.target_url,
    l.shortened_url,
    l.status,
    l.expiration_date,
    l.created_at,
    l.updated_at
FROM links l
JOIN users u ON l.user_id = u.id
WHERE status = $1 AND deleted_at IS NULL;

-- name: UpdateLink :one
UPDATE links l
SET
    target_url = $2,
    status = $3,
    expiration_date = $4,
    updated_at = $5
FROM users u
WHERE l.id = $1
  AND l.user_id = u.id
  AND l.delete_at IS NULL
RETURNING
    l.id,
    u.user_name,
    l.target_url,
    l.shortened_url,
    l.status,
    l.expiration_date,
    l.created_at,
    l.updated_at;

-- name: UpdateLinkShortenedURL :one
UPDATE links l
SET
    shortened_url = $2
FROM users u
WHERE l.id = $1
    AND l.user_id = u.id
    AND l.deleted_at IS NULL
RETURNING
    l.id,
    u.user_name,
    l.target_url,
    l.shortened_url,
    l.status,
    l.expiration_date,
    l.created_at,
    l.updated_at;

-- name: DeleteLink :exec
UPDATE links
SET
    deleted_at = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreLink :exec
UPDATE links
SET
    deleted_at = NULL
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: DestroyLink :exec
DELETE FROM links
WHERE id = $1;
