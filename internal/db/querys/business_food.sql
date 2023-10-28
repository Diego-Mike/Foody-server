-- name: CreateNewFood :one
INSERT INTO business_food (business_id, food_img, food_title, food_description, food_price, food_available_per_day) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
