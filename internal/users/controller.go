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
	currentUser := r.Context().Value(constants.UserContextKey).(mw.JwtUserData)

	resp := config.ClientResponse{Rsp: struct {
		User mw.JwtUserData `json:"user"`
	}{
		User: currentUser,
	}}
	config.WriteResponse(w, http.StatusOK, resp)
}
