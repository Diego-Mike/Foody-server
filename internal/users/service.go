package users

import (
	"net/http"

	"github.com/Foody-App-Tech/Main-server/config"
	"github.com/Foody-App-Tech/Main-server/internal/constants"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
	mw "github.com/Foody-App-Tech/Main-server/internal/global_middlewares"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

func (service *UserService) getFullReservation(r *http.Request) (userReservation, string, error) {

	user := r.Context().Value(constants.UserContextKey).(mw.JwtUserData)

	reservation, err := service.storage.GetUserReservation(r.Context(), user.UserID)
	if err != nil {
		return userReservation{}, "Problemas obteniendo la información !", err
	}

	prettifyCash := message.NewPrinter(language.English)
	var reservationFoodRestructured userReservation
	for _, v := range reservation {
		if len(reservationFoodRestructured.Foods) == 0 {
			reservationFoodRestructured.ReservationID = v.ReservationID
			reservationFoodRestructured.BusinessID = v.BusinessID
			reservationFoodRestructured.OrderSchedule = v.OrderSchedule.Time
			reservationFoodRestructured.CreatedAt = v.CreatedAt.Time
		}

		prettyCash := prettifyCash.Sprintf("%d", v.FoodPrice)
		reservationFoodRestructured.Foods = append(reservationFoodRestructured.Foods, reservedFood{
			FoodID:    v.FoodID,
			FoodTitle: v.FoodTitle,
			FoodPrice: foodPrice{
				Prettify:  prettyCash,
				RealPrice: v.FoodPrice,
			},
			FoodImg:         v.FoodImg,
			Amount:          v.Amount,
			FoodDetails:     v.Details.String,
			FoodDescription: v.FoodDescription.String,
		})
	}

	return reservationFoodRestructured, "", nil

}
