// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	AddBusinessMember(ctx context.Context, arg AddBusinessMemberParams) (int64, error)
	AddFoodsToReservation(ctx context.Context, arg AddFoodsToReservationParams) (int64, error)
	CreateBusiness(ctx context.Context, arg CreateBusinessParams) (Business, error)
	CreateNewFood(ctx context.Context, arg CreateNewFoodParams) (BusinessFood, error)
	CreateNewNotification(ctx context.Context, arg CreateNewNotificationParams) (BusinessReservationsNotificacion, error)
	// Creating new order
	CreateReservation(ctx context.Context, arg CreateReservationParams) (BusinessReservation, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GenerateSession(ctx context.Context, arg GenerateSessionParams) (Session, error)
	GetBusinessById(ctx context.Context, businessID int64) (Business, error)
	GetFullUser(ctx context.Context, userID int64) (GetFullUserRow, error)
	GetHomeBusinessFood(ctx context.Context, arg GetHomeBusinessFoodParams) ([]GetHomeBusinessFoodRow, error)
	GetNextHomePage(ctx context.Context, businessID int64) (sql.NullInt64, error)
	GetSession(ctx context.Context, userIDSession int64) (Session, error)
	GetUserById(ctx context.Context, userID int64) (User, error)
	GetUserBySocialId(ctx context.Context, socialID string) (User, error)
	// Getting order
	GetUserReservation(ctx context.Context, userID int64) ([]GetUserReservationRow, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
