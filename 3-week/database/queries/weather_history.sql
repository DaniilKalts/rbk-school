-- name: CreateWeatherHistory :one
INSERT INTO weather_history (id, user_id, city, temperature, description)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, city, temperature, description, requested_at;
