package businesses

import (
	"github.com/Foody-App-Tech/Main-server/config"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
)

type BusinessesService struct {
	storage db.Store
	env     config.EnvVariables
}

func NewBusinessesService(storage db.Store, env config.EnvVariables) *BusinessesService {
	return &BusinessesService{
		storage: storage,
		env:     env,
	}
}
