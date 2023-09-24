package users

import (
	"net/http"
)

type loginReq struct {
	ACCESS_TOKEN string `json:"access_token"`
}

func checkingLoginPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// // read req payload
		// var reqPayload loginReq
		// err := config.ReadBody(w, r, &reqPayload)
		// if err != nil {
		// 	config.ErrorResponse(w, fmt.Sprintf("Error reading payload --> %s", err), http.StatusBadRequest)
		// 	return
		// }

		// // validate data
		// validationErr := validation.ValidateStruct(&reqPayload,
		// 	validation.Field(&reqPayload.Email, validation.Required, is.Alphanumeric, is.Email),
		// 	validation.Field(&reqPayload.Provider, validation.Required, is.Alphanumeric, validation.In("Google", "Twitter", "Facebook")),
		// 	validation.Field(&reqPayload.Username, validation.Required, is.Alphanumeric))

		// if validationErr != nil {
		// 	config.ErrorResponse(w, fmt.Sprintf("Error validating payload --> %s", validationErr), http.StatusBadRequest)
		// 	return
		// }
	})

}
