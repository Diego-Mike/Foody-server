-- name: CreateBusiness :one
INSERT INTO businesses (name, city, address, latitude, longitude, ubication_photo) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: AddBusinessMember :one
INSERT INTO business_members (business_id, user_id, business_position) VALUES ($1, $2, $3) 
RETURNING *;