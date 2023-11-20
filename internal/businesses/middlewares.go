package businesses

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Foody-App-Tech/Main-server/config"
	"github.com/Foody-App-Tech/Main-server/internal/constants"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var queryDecoder = schema.NewDecoder()

type getBusinessHomeFoodRequest struct {
	PageSize      int64 `validate:"gt=0,required"`
	AfterBusiness int64 `validate:"gt=-1,numeric"`
}

func checkGetBusinessHomeFoodPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// payload
		var reqPayload getBusinessHomeFoodRequest
		err := queryDecoder.Decode(&reqPayload, r.URL.Query())
		if err != nil {
			config.ErrorResponse(w, "Hubo un problema leyendo los parametros de la petición", err, http.StatusBadRequest)
			return
		}

		payloadValidationErr := config.ValidateData(reqPayload)
		if payloadValidationErr != nil {
			config.ErrorResponse(w, "Hubo un problema leyendo los parametros de la petición", payloadValidationErr, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), constants.RequestPayloadKey, reqPayload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type createNewBusinessRequest struct {
	Name             string  `json:"name" validate:"required,max=30"`
	City             string  `json:"city" validate:"required,max=50"`
	Address          string  `json:"address" validate:"required,max=100"`
	Latitude         string  `json:"latitude" validate:"required,max=100"`
	Longitude        string  `json:"longitude" validate:"required,max=100"`
	Presentation     string  `json:"presentation" validate:"required,max=300"`
	ClientsMaxAmount int16   `json:"clients_max_amount" validate:"omitempty,min=0"`
	DaysOfWeek       []int16 `json:"days_of_week" validate:"dive,required,numeric,min=1,max=7"`
	UserID           int64   `json:"user_id" validate:"required,min=1"`
	BusinessPosition string  `json:"business_position" validate:"required,oneof=Dueño Administrador Empleado"`
}

func checkNewBusinessPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// read payload
		var reqPayload createNewBusinessRequest
		err := config.ReadBody(w, r, &reqPayload)
		if err != nil {
			config.ErrorResponse(w, fmt.Sprintf("error reading payload ----> %s", err), nil, http.StatusBadRequest)
			return
		}

		// validate payload
		payloadValidationErr := config.ValidateData(reqPayload)
		if payloadValidationErr != nil {
			config.ErrorResponse(w, "there where problems with the payload", payloadValidationErr, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), constants.RequestPayloadKey, reqPayload)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

type createNewBusinessMemberPayload struct {
	UserID           int64  `json:"user_id" validate:"required,min=1"`
	BusinessPosition string `json:"business_position" validate:"required,oneof=Dueño Administrador Empleado"`
}

type createNewBusinessMemberRequest struct {
	businessIdParameter
	createNewBusinessMemberPayload
}

func checkNewBusinessMemberPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// FIXME: unify validations

		// get param
		getBusinessId := chi.URLParam(r, "business-id")
		if getBusinessId == "" {
			config.ErrorResponse(w, "error getting business id parameter", nil, http.StatusBadRequest)
			return
		}
		businessId, err := strconv.ParseInt(getBusinessId, 10, 64)
		if err != nil {
			config.ErrorResponse(w, fmt.Sprintf("error converting business id parameter to int: %s", err), nil, http.StatusBadRequest)
			return
		}

		// validate param
		paramToStruct := businessIdParameter{BusinessID: businessId}
		paramValidationErr := config.ValidateData(paramToStruct)
		if paramValidationErr != nil {
			config.ErrorResponse(w, "there where problems with the parameter", paramValidationErr, http.StatusBadRequest)
			return
		}

		// read payload
		var reqPayload createNewBusinessMemberPayload
		err = config.ReadBody(w, r, &reqPayload)
		if err != nil {
			config.ErrorResponse(w, fmt.Sprintf("error reading payload: %s", err), nil, http.StatusBadRequest)
			return
		}

		// validate payload
		payloadValidationErr := config.ValidateData(reqPayload)
		if payloadValidationErr != nil {
			config.ErrorResponse(w, "payload validation errors", payloadValidationErr, http.StatusBadRequest)
			return
		}
		fullPayload := createNewBusinessMemberRequest{businessIdParameter: paramToStruct, createNewBusinessMemberPayload: reqPayload}

		ctx := context.WithValue(r.Context(), constants.RequestPayloadKey, fullPayload)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func checkGetBusinessByIdPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get param
		getBusinessId := chi.URLParam(r, "business-id")
		if getBusinessId == "" {
			config.ErrorResponse(w, "error getting business id parameter", nil, http.StatusBadRequest)
			return
		}
		businessId, err := strconv.ParseInt(getBusinessId, 10, 64)
		if err != nil {
			config.ErrorResponse(w, fmt.Sprintf("business-id has invalid type: %s", err), nil, http.StatusBadRequest)
			return
		}

		// validate param
		paramToStruct := businessIdParameter{BusinessID: businessId}
		paramValidationErr := config.ValidateData(paramToStruct)
		if paramValidationErr != nil {
			config.ErrorResponse(w, "there where problems with param validation", paramValidationErr, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), constants.RequestPayloadKey, paramToStruct)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

type createFoodPayload struct {
	FoodImg             string `json:"food_img" validate:"required"`
	FoodTitle           string `json:"food_title" validate:"required,max=55"`
	FoodDescription     string `json:"food_description" validate:"omitempty,max=150"`
	FoodPrice           int64  `json:"food_price" validate:"required,gt=0"`
	FoodAvailablePerDay int16  `json:"food_available_per_day" validate:"omitempty,gt=0"`
}

type createFoodRequest struct {
	businessIdParameter
	createFoodPayload
}

func checkCreateFoodPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		businessId, err := getBusinessId(r)
		if err != nil {
			config.ErrorResponse(w, err.Error(), nil, http.StatusBadRequest)
			return
		}

		// read payload
		var reqPayload createFoodPayload
		err = config.ReadBody(w, r, &reqPayload)
		if err != nil {
			config.ErrorResponse(w, fmt.Sprintf("error reading payload: %s", err), nil, http.StatusBadRequest)
			return
		}

		// validate payload
		fullPayload := createFoodRequest{businessIdParameter: businessId, createFoodPayload: reqPayload}
		payloadValidationErr := config.ValidateData(fullPayload)
		if payloadValidationErr != nil {
			config.ErrorResponse(w, "payload validation errors", payloadValidationErr, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), constants.RequestPayloadKey, fullPayload)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
