// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"database/sql"
	"time"
)

type Business struct {
	BusinessID int64  `json:"business_id"`
	Name       string `json:"name"`
	City       string `json:"city"`
	Address    string `json:"address"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	// max amount of people that can order food in this business
	Presentation string `json:"presentation"`
	// number of people that can be inside the business, field is not mandatory
	ClientsMaxAmount sql.NullInt16 `json:"clients_max_amount"`
	CreatedAt        time.Time     `json:"created_at"`
}

type BusinessFood struct {
	BusinessID      int64          `json:"business_id"`
	FoodID          int64          `json:"food_id"`
	FoodImg         string         `json:"food_img"`
	FoodTitle       string         `json:"food_title"`
	FoodDescription sql.NullString `json:"food_description"`
	FoodPrice       int64          `json:"food_price"`
	// number of pieces of this food that can be sold per day, it is not mandatory
	FoodAvailablePerDay sql.NullInt16 `json:"food_available_per_day"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}

// primary key 🔑 is composed by business_id and user_id
type BusinessMember struct {
	BusinessID int64 `json:"business_id"`
	UserID     int64 `json:"user_id"`
	// posible positions for the employee, dueño, administrador, empleado
	BusinessPosition string    `json:"business_position"`
	CreatedAt        time.Time `json:"created_at"`
}

type BusinessReservation struct {
	ReservationID int64        `json:"reservation_id"`
	BusinessID    int64        `json:"business_id"`
	UserID        int64        `json:"user_id"`
	OrderSchedule sql.NullTime `json:"order_schedule"`
	Accepted      bool         `json:"accepted"`
	CreatedAt     sql.NullTime `json:"created_at"`
}

type BusinessReservationsNotificacion struct {
	ReservationID           int64        `json:"reservation_id"`
	NotificationTitle       string       `json:"notification_title"`
	NotificationDescription string       `json:"notification_description"`
	CreatedAt               sql.NullTime `json:"created_at"`
}

type BusinessReservationsState struct {
	ReservationID         int64        `json:"reservation_id"`
	CancelledByClient     bool         `json:"cancelled_by_client"`
	CancelledByBusiness   bool         `json:"cancelled_by_business"`
	ReasonForCancellation string       `json:"reason_for_cancellation"`
	CreatedAt             sql.NullTime `json:"created_at"`
}

type ReserveFood struct {
	ReservationID int64          `json:"reservation_id"`
	FoodID        int64          `json:"food_id"`
	Amount        int16          `json:"amount"`
	Details       sql.NullString `json:"details"`
	CreatedAt     sql.NullTime   `json:"created_at"`
}

type Session struct {
	UserIDSession int64     `json:"user_id_session"`
	Valid         bool      `json:"valid"`
	UserAgent     string    `json:"user_agent"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// primary key is composed by user_id and social_id, both of those properties are unique
type User struct {
	UserID int64 `json:"user_id"`
	// this field is for the unique id provided by google, or facebook, or tik tok etc etc.. to identify a user
	SocialID     string    `json:"social_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Picture      string    `json:"picture"`
	Provider     string    `json:"provider"`
	RegisteredAt time.Time `json:"registered_at"`
}
