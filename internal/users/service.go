package users

import (
	"github.com/Foody-App-Tech/Main-server/config"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
)

type UserService struct {
	storage *db.SQLCStore
	env     config.EnvVariables
}

func NewUserService(storage *db.SQLCStore, env config.EnvVariables) *UserService {
	return &UserService{
		storage: storage,
		env:     env,
	}
}
