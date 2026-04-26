-- name: CreateWeatherHistory :one
INSERT INTO weather_history (id, user_id, city, temperature, description)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, city, temperature, description, requested_at;

-- name: ListWeatherHistoryByUserAndCity :many
SELECT id, user_id, city, temperature, description, requested_at
FROM weather_history
WHERE user_id = $1 AND city = $2
ORDER BY requested_at DESC
LIMIT NULLIF($3, 0);
