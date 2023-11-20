-- name: CreateBusiness :one
INSERT INTO businesses (name, city, address, latitude, longitude, presentation, clients_max_amount) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: AddBusinessMember :one
INSERT INTO business_members (business_id, user_id, business_position) VALUES ($1, $2, $3) 
RETURNING user_id;

-- name: GetBusinessById :one
SELECT * FROM businesses WHERE business_id = $1;

-- name: GetHomeBusinessFood :many
SELECT b.business_id, b."name", b.city, bf.food_id, bf.food_title, bf.food_description, bf.food_price, bf.food_available_per_day, 
bf.food_img FROM businesses b 
INNER JOIN lateral (
    SELECT bf.food_id, bf.food_title, bf.food_description, bf.food_price, bf.food_available_per_day, bf.food_img FROM business_food bf 
    where b.business_id = bf.business_id 
    ORDER BY bf.created_at DESC
    LIMIT 3
) bf ON true
WHERE b.business_id >= sqlc.arg(after_business)::bigint
LIMIT sqlc.arg(page_size)::bigint;

-- name: GetNextHomePage :one
SELECT b.business_id from businesses b 
RIGHT JOIN business_food bf ON b.business_id = bf.business_id
WHERE b.business_id > $1 
LIMIT 1;
