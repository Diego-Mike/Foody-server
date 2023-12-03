-- Creating new order

-- name: CreateReservation :one
INSERT INTO business_reservations (business_id, user_id, order_schedule) VALUES($1, $2, $3)
RETURNING *;

-- name: AddFoodsToReservation :execrows
INSERT INTO reserve_food (reservation_id, food_id, amount, details) VALUES($1, $2, $3, $4) RETURNING *;

-- name: CreateNewNotification :one
INSERT INTO business_reservations_notificacions (reservation_id, notification_title, notification_description) VALUES($1, $2, $3)
RETURNING *;

-- Getting order

-- name: GetUserReservation :many
SELECT br.reservation_id, br.business_id, br.created_at, br.order_schedule, rf.food_id,
bf.food_title,  bf.food_price, bf.food_img, rf.amount, rf.details, bf.food_description
FROM business_reservations br 
INNER JOIN reserve_food rf ON br.reservation_id = rf.reservation_id 
INNER JOIN business_food bf ON rf.food_id = bf.food_id
WHERE br.user_id = $1 and br.created_at::date = CURRENT_DATE;

-- Adding description to the order

-- Accepting order from business

-- Cancelling order from business

-- Cancelling order from user
