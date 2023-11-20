package businesses

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Foody-App-Tech/Main-server/config"
	"github.com/Foody-App-Tech/Main-server/internal/constants"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
	mw "github.com/Foody-App-Tech/Main-server/internal/global_middlewares"
)

type BusinessesController struct {
	businessesService *BusinessesService
	globalHelpers     *mw.GlobalMiddlewares
}

func NewBusinessesController(businessesService *BusinessesService, globalHelpers *mw.GlobalMiddlewares) *BusinessesController {
	return &BusinessesController{
		businessesService: businessesService,
		globalHelpers:     globalHelpers,
	}
}

type Food struct {
	FoodID              int64  `json:"food_id"`
	FoodTitle           string `json:"food_title"`
	FoodDescription     string `json:"food_description"`
	FoodPrice           string `json:"food_price"`
	FoodAvailablePerDay int16  `json:"food_available_per_day"`
	FoodImg             string `json:"food_img"`
}

type businessHomeFood struct {
	BusinessID int64  `json:"business_id"`
	Name       string `json:"name"`
	City       string `json:"city"`
	Foods      []Food `json:"foods"`
}

type homeFoodRsp struct {
	HomeFood []businessHomeFood `json:"home_food"`
	NextPage int64              `json:"next_page"`
}

func (c *BusinessesController) getBusinessHomeFood(w http.ResponseWriter, r *http.Request) {

	homeFoodRestructured, msgErr, dbErr := c.businessesService.getFood(r)
	if dbErr != nil {
		config.ErrorResponse(w, msgErr, dbErr, http.StatusServiceUnavailable)
		return
	}
	// log.Println("home food re structure", config.PrettyPrint(homeFoodRestructured))

	var lastBusinessId int64
	if len(homeFoodRestructured) > 0 {
		lastBusinessId = homeFoodRestructured[len(homeFoodRestructured)-1].BusinessID
	}
	nextPage, err := c.businessesService.storage.GetNextHomePage(r.Context(), lastBusinessId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("problem gettin next page for home", err)
		config.ErrorResponse(w, "Problemas al obtener la siguiente pagina !", err.Error(), http.StatusServiceUnavailable)
		return
	}

	rsp := config.ClientResponse{
		Rsp: homeFoodRsp{HomeFood: homeFoodRestructured, NextPage: nextPage.Int64},
	}
	config.WriteResponse(w, http.StatusOK, rsp)

}

// FIXME: should we test this
func (c *BusinessesController) newBusiness(w http.ResponseWriter, r *http.Request) {

	createdBusiness, err := c.businessesService.createNewBusiness(r)
	if err != nil {
		config.ErrorResponse(w, err.Error(), nil, http.StatusServiceUnavailable)
		return
	}

	rsp := config.ClientResponse{Rsp: db.CreateNewBusinessTxResult{
		BusinessId: createdBusiness.BusinessId,
	}}
	config.WriteResponse(w, http.StatusCreated, rsp)

}

type newMemberRsp struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}

func (c *BusinessesController) newMember(w http.ResponseWriter, r *http.Request) {

	memberId, statusCode, err := c.businessesService.createNewBusinessMember(r)
	if err != nil {
		config.ErrorResponse(w, err.Error(), nil, statusCode)
		return
	}

	user, err := c.businessesService.storage.GetUserById(r.Context(), memberId)
	if err != nil {
		log.Println("problem getting user:", err)
		config.ErrorResponse(w, fmt.Sprintf("ocurrio un problema al intentar obtener el nuevo miembro del negocio: %s", err), nil, http.StatusInternalServerError)
		return
	}

	rsp := config.ClientResponse{Rsp: newMemberRsp{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Picture:  user.Picture,
	}}

	config.WriteResponse(w, http.StatusCreated, rsp)

}

type getBusinessRsp struct {
	Name      string `json:"name"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	// max amount of people that can order food in this business
	Presentation     string `json:"presentation"`
	ClientsMaxAmount int16  `json:"clients_max_amount"`
}

func (c *BusinessesController) getBusinessById(w http.ResponseWriter, r *http.Request) {

	businessId := r.Context().Value(constants.RequestPayloadKey).(businessIdParameter)

	business, err := c.businessesService.storage.GetBusinessById(r.Context(), businessId.BusinessID)
	if err != nil {
		// business does not exist
		if err == sql.ErrNoRows {
			log.Printf("business %s does not exist: %s", strconv.Itoa(int(businessId.BusinessID)), err)
			config.ErrorResponse(w, "Este negocio no existe en nuestro sistema", nil, http.StatusNotFound)
			return
		}
		// there was a problem with the db
		log.Printf("there was a problem gettin business %s: %s", strconv.Itoa(int(businessId.BusinessID)), err)
		config.ErrorResponse(w, "Hay problemas con nuestros servidores", nil, http.StatusServiceUnavailable)
		return
	}

	rsp := config.ClientResponse{Rsp: getBusinessRsp{Name: business.Name, City: business.City, Address: business.Address, Latitude: business.Latitude, Longitude: business.Longitude, Presentation: business.Presentation, ClientsMaxAmount: business.ClientsMaxAmount.Int16}}
	config.WriteResponse(w, http.StatusOK, rsp)

}

type newFoodRsp struct {
	FoodID              int64  `json:"food_id"`
	FoodImg             string `json:"food_img"`
	FoodTitle           string `json:"food_title"`
	FoodDescription     string `json:"food_description"`
	FoodPrice           string `json:"food_price"`
	FoodAvailablePerDay int16  `json:"food_available_per_day"`
}

func (c *BusinessesController) createFood(w http.ResponseWriter, r *http.Request) {

	// TODO: check if member is allowed to do this
	newFood, statusCode, err := c.businessesService.createNewBusinessFood(r)
	if err != nil {
		config.ErrorResponse(w, err.Error(), nil, statusCode)
		return
	}

	config.WriteResponse(w, http.StatusCreated, newFood)

}
