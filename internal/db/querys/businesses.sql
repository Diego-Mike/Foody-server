-- name: CreateBusiness :one
INSERT INTO businesses (name, city, address, latitude, longitude, presentation, clients_max_amount) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: AddBusinessSchedule :execrows
INSERT INTO business_schedule (day_of_week, opening_hour, closing_hour) VALUES (UNNEST(@days_of_week::smallint[]), UNNEST(@opening_hours::time[]), UNNEST(@closing_hours::time[]));

-- name: AddBusinessMember :one
INSERT INTO business_members (business_id, user_id, business_position) VALUES ($1, $2, $3) 
RETURNING user_id;

-- name: GetBusinessById :one
SELECT * FROM businesses WHERE business_id = $1;
