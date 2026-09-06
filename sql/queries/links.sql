-- name: CreateLink :one
INSERT INTO links (user_id, target_url, shortened_url, password, status, expiration_date, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetLinks :many
SELECT
    id,
    user_id,
    target_url,
    shortened_url,
    status,
    expiration_date,
    created_at,
    updated_at
FROM links
WHERE deleted_at IS NULL;

-- name: GetLinkByTargetURL :one
SELECT
    id,
    user_id,
    target_url,
    shortened_url,
    status,
    expiration_date,
    created_at,
    updated_at
FROM links
WHERE target_url = $1 AND deleted_at IS NULL;

-- name: GetLinkByShortenedURL :one
SELECT
    id,
    user_id,
    target_url,
    shortened_url,
    status,
    expiration_date,
    created_at,
    updated_at
FROM links
WHERE shortened_url = $1 AND deleted_at IS NULL;

-- name: GetLinksByStatus :many
SELECT
    id,
    user_id,
    target_url,
    shortened_url,
    status,
    expiration_date,
    created_at,
    updated_at
FROM links
WHERE status = $1 AND deleted_at IS NULL;

-- name: UpdateLink :one
UPDATE links
SET
    target_url = $2,
    status = $3,
    expiration_date = $4,
    updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteLink :exec
UPDATE links
SET
    deleted_at = $2
WHERE id = $1;

-- name: RestoreLink :exec
UPDATE links
SET
    deleted_at = NULL
WHERE id = $1;

-- name: DestroyLink :exec
DELETE FROM links
WHERE id = $1;

-- name: NextLinkID :one
SELECT nextval('links_id_seq');
