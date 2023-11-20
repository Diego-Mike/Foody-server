package businesses

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Foody-App-Tech/Main-server/config"
	"github.com/Foody-App-Tech/Main-server/internal/constants"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
	mw "github.com/Foody-App-Tech/Main-server/internal/global_middlewares"
	"github.com/lib/pq"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

// TODO: test services functions

func (service *BusinessesService) getFood(r *http.Request) ([]businessHomeFood, string, error) {
	payload := r.Context().Value(constants.RequestPayloadKey).(getBusinessHomeFoodRequest)

	arg := db.GetHomeBusinessFoodParams{
		AfterBusiness: payload.AfterBusiness,
		PageSize:      payload.PageSize,
	}

	homeFood, err := service.storage.GetHomeBusinessFood(r.Context(), arg)
	if err != nil {
		return nil, "Problemas obteniendo la información !", err
	}

	prettifyCash := message.NewPrinter(language.English)
	var homeFoodRestructured []businessHomeFood
	for _, v := range homeFood {
		prettyCash := prettifyCash.Sprintf("%d", v.FoodPrice)
		food := Food{
			FoodID: v.FoodID, FoodTitle: v.FoodTitle, FoodDescription: v.FoodDescription.String, FoodPrice: prettyCash, FoodAvailablePerDay: v.FoodAvailablePerDay.Int16, FoodImg: v.FoodImg,
		}

		if len(homeFoodRestructured) == 0 {
			homeFoodRestructured = append(homeFoodRestructured, businessHomeFood{BusinessID: v.BusinessID, Name: v.Name, City: v.City,
				Foods: []Food{food}})
			continue
		}

		for i, v2 := range homeFoodRestructured {
			// add food to business
			if v2.BusinessID == v.BusinessID {
				homeFoodRestructured[i].Foods = append(homeFoodRestructured[i].Foods, food)
				break
			} else if len(homeFoodRestructured) == (i + 1) {
				// add new business + first food
				homeFoodRestructured = append(homeFoodRestructured, businessHomeFood{BusinessID: v.BusinessID, Name: v.Name, City: v.City,
					Foods: []Food{food}})
				break
			} else {
				continue

			}
		}
	}

	return homeFoodRestructured, "", nil

}

func (service *BusinessesService) createNewBusiness(r *http.Request) (db.CreateNewBusinessTxResult, error) {
	userReq := r.Context().Value(constants.UserContextKey).(mw.JwtUserData)
	payload := r.Context().Value(constants.RequestPayloadKey).(createNewBusinessRequest)

	arg := db.CreateNewBusinessTxParams{
		CreateBusinessParams: db.CreateBusinessParams{
			Name:         payload.Name,
			City:         payload.City,
			Address:      payload.Address,
			Latitude:     payload.Latitude,
			Longitude:    payload.Longitude,
			Presentation: payload.Presentation,
			ClientsMaxAmount: sql.NullInt16{
				Int16: payload.ClientsMaxAmount,
				Valid: true,
			},
		},
		UserID:           userReq.UserID,
		BusinessPosition: payload.BusinessPosition,
	}

	createdBusiness, err := service.storage.CreateNewBusinessTx(r.Context(), arg)
	if err != nil {
		log.Println("problem creating new business:", err)

		err = fmt.Errorf("ocurrio un problema creando el negocio: %s", err)
		return db.CreateNewBusinessTxResult{}, err
	}

	return createdBusiness, nil
}

func (service *BusinessesService) createNewBusinessMember(r *http.Request) (int64, int, error) {

	userReq := r.Context().Value(constants.UserContextKey).(mw.JwtUserData)
	payload := r.Context().Value(constants.RequestPayloadKey).(createNewBusinessMemberRequest)

	arg := db.AddBusinessMemberParams{
		BusinessID:       payload.BusinessID,
		UserID:           userReq.UserID,
		BusinessPosition: payload.BusinessPosition,
	}

	businessMember, err := service.storage.AddBusinessMember(r.Context(), arg)
	if err != nil {
		log.Println("problem creating new business member:", err)

		// business member already exists
		if prErr, ok := err.(*pq.Error); ok && prErr.Code == "23505" {
			err = fmt.Errorf("el miembro ya existe en el negocio: %s", err)
			return 0, http.StatusBadRequest, err
		}

		err = fmt.Errorf("ocurrio un problema al intentar agregar un miembro al negocio: %s", err)
		return 0, http.StatusServiceUnavailable, err
	}

	return businessMember, 0, nil
}

func (service *BusinessesService) createNewBusinessFood(r *http.Request) (newFoodRsp, int, error) {
	payload := r.Context().Value(constants.RequestPayloadKey).(createFoodRequest)

	arg := db.CreateNewFoodParams{BusinessID: payload.BusinessID,
		FoodImg:             payload.FoodImg,
		FoodTitle:           payload.FoodTitle,
		FoodDescription:     sql.NullString{String: payload.FoodDescription, Valid: true},
		FoodPrice:           payload.FoodPrice,
		FoodAvailablePerDay: sql.NullInt16{Int16: payload.FoodAvailablePerDay, Valid: true}}

	newFood, dbErr := service.storage.CreateNewFood(r.Context(), arg)
	if dbErr != nil {

		log.Println("problem creating new food:", dbErr)

		if pqErr, ok := dbErr.(*pq.Error); ok && pqErr.Code == "23503" {
			// business does not exist
			err := fmt.Errorf("el negocio al cuál se le quiere crear la nueva comida no existe : %s", dbErr)
			return newFoodRsp{}, http.StatusBadRequest, err
		}
		// a problem with the db
		err := fmt.Errorf("problemas al crear la nueva comida : %s", dbErr)
		return newFoodRsp{}, http.StatusServiceUnavailable, err
	}

	// log.Println("food created", config.PrettyPrint(newFood))

	prettifyCash := message.NewPrinter(language.English)
	prettyCash := prettifyCash.Sprintf("%d", newFood.FoodPrice)

	return newFoodRsp{
		FoodID:              newFood.FoodID,
		FoodImg:             newFood.FoodImg,
		FoodTitle:           newFood.FoodTitle,
		FoodDescription:     newFood.FoodDescription.String,
		FoodPrice:           prettyCash,
		FoodAvailablePerDay: newFood.FoodAvailablePerDay.Int16,
	}, 0, nil
}
