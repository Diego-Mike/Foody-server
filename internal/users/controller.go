package users

import (
	"net/http"

	"github.com/Foody-App-Tech/Main-server/config"
	"github.com/Foody-App-Tech/Main-server/internal/constants"
	mw "github.com/Foody-App-Tech/Main-server/internal/global_middlewares"
)

type UserController struct {
	userService   *UserService
	globalHelpers *mw.GlobalMiddlewares
}

func NewUserController(userService *UserService, globalHelpers *mw.GlobalMiddlewares) *UserController {
	return &UserController{userService: userService, globalHelpers: globalHelpers}
}

func (u *UserController) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(constants.UserContextKey)
	if currentUser == nil {
		config.ErrorResponse(w, "unauthorized to get user data", http.StatusUnauthorized)
		return
	}

	resp := config.ClientResponse{Error: false, Message: "successfull call", Data: currentUser}
	config.WriteResponse(w, http.StatusOK, resp)
}
