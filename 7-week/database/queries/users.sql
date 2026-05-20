-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email, password_hash, salt, role)
VALUES ($1, $2, $3, $4, $5, $6,$7)
RETURNING id, first_name, last_name, email, role, created_at, updated_at, deleted_at;

-- name: GetUserByID :one
SELECT id,
       first_name,
       last_name,
       email,
       role,
       created_at,
       updated_at,
       deleted_at
FROM users
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT id,
       first_name,
       last_name,
       email,
       role,
       created_at,
       updated_at,
       deleted_at
FROM users
WHERE email = $1
  AND deleted_at IS NULL;

-- name: GetUserCredentialsByEmail :one
SELECT id,
       email,
       password_hash,
       salt,
       role
FROM users
WHERE email = $1
  AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id,
       first_name,
       last_name,
       email,
       role,
       created_at,
       updated_at,
       deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2,
    last_name  = $3,
    email      = $4,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, first_name, last_name, email, role, created_at, updated_at, deleted_at;

-- name: SoftDeleteUser :execrows
UPDATE users
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;
