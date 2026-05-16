-- name: CreateUserCity :one
INSERT INTO user_cities (id, user_id, city)
VALUES ($1, $2, $3)
RETURNING id, user_id, city, created_at;

-- name: ListUserCities :many
SELECT id, user_id, city, created_at
FROM user_cities
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteUserCity :execrows
DELETE
FROM user_cities
WHERE id = $1
  AND user_id = $2;
