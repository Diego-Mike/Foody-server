package businesses

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type businessIdParameter struct {
	BusinessID int64 `json:"business_id" validate:"required,min=1"`
}

func getBusinessId(r *http.Request) (businessIdParameter, error) {
	getBusinessId := chi.URLParam(r, "business-id")
	if getBusinessId == "" {
		err := errors.New("error getting business id parameter")
		return businessIdParameter{}, err
	}
	businessId, err := strconv.ParseInt(getBusinessId, 10, 64)
	if err != nil {
		err := fmt.Errorf("business id has invalid type: %s", err)
		return businessIdParameter{}, err
	}
	paramToStruct := businessIdParameter{BusinessID: businessId}

	return paramToStruct, nil
}

type reservationIdParameter struct {
	ReservationID int64 `json:"reservation_id" validate:"required,min=1"`
}

func getReservationId(r *http.Request) (reservationIdParameter, error) {
	getReservationId := chi.URLParam(r, "reservation-id")

	if getReservationId == "" {
		return reservationIdParameter{}, errors.New("id de la reservacion no es correcto")
	}
	reservationId, err := strconv.ParseInt(getReservationId, 10, 64)
	if err != nil {
		return reservationIdParameter{}, errors.New("id de la reservación no tiene type correcto")
	}

	paramToStruct := reservationIdParameter{ReservationID: reservationId}

	return paramToStruct, nil

}
