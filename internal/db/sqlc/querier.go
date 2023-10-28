// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"context"
)

type Querier interface {
	AddBusinessMember(ctx context.Context, arg AddBusinessMemberParams) (int64, error)
	AddBusinessSchedule(ctx context.Context, arg AddBusinessScheduleParams) (int64, error)
	CreateBusiness(ctx context.Context, arg CreateBusinessParams) (Business, error)
	CreateNewFood(ctx context.Context, arg CreateNewFoodParams) (BusinessFood, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GenerateSession(ctx context.Context, arg GenerateSessionParams) (Session, error)
	GetBusinessById(ctx context.Context, businessID int64) (Business, error)
	GetSession(ctx context.Context, userIDSession int64) (Session, error)
	GetUserById(ctx context.Context, userID int64) (User, error)
	GetUserBySocialId(ctx context.Context, socialID string) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
