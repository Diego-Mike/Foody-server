package users

import (
	"database/sql"
	"fmt"
	"log"
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

type FullUserRsp struct {
	UserID           int64  `json:"user_id"`
	SocialID         string `json:"social_id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Picture          string `json:"picture"`
	BusinessID       int64  `json:"business_id,omitempty"`
	BusinessPosition string `json:"business_position,omitempty"`
	IsBusinessMember bool   `json:"is_business_member"`
}

func (u *UserController) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(constants.UserContextKey).(mw.JwtUserData)

	fullUser, err := u.userService.storage.GetFullUser(r.Context(), currentUser.UserID)
	if err != nil {
		log.Println("Error getting full user", err)
		// user does not exist
		if err == sql.ErrNoRows {
			config.ErrorResponse(w, "El usuario no existe", err, http.StatusServiceUnavailable)
			return
		}
		// db problem
		config.ErrorResponse(w, fmt.Sprintf("Ocurrio un problema al intentar obtener el usuario: %s", err), nil, http.StatusServiceUnavailable)
		return
	}

	var user FullUserRsp
	user.Username = fullUser.Username
	user.UserID = fullUser.UserID
	user.SocialID = fullUser.SocialID
	user.Picture = fullUser.Picture
	user.Email = fullUser.Email

	// business information from the user
	if fullUser.BusinessPosition.Valid && fullUser.BusinessID.Valid {
		user.BusinessPosition = fullUser.BusinessPosition.String
		user.BusinessID = fullUser.BusinessID.Int64
		user.IsBusinessMember = true
	}

	resp := config.ClientResponse{Rsp: struct {
		User FullUserRsp `json:"user"`
	}{
		User: user,
	}}

	config.WriteResponse(w, http.StatusOK, resp)
}
