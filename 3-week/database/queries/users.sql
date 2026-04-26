-- name: CreateUser :one
INSERT INTO users (id, name, email)
VALUES ($1, $2, $3)
RETURNING id, name, email, created_at, updated_at, deleted_at;

-- name: GetUserByID :one
SELECT id, name, email, created_at, updated_at, deleted_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id, name, email, created_at, updated_at, deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET name = $2,
    email = $3,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, name, email, created_at, updated_at, deleted_at;

-- name: SoftDeleteUser :execrows
UPDATE users
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
