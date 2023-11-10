package login

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Foody-App-Tech/Main-server/config"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
)

type LoginService struct {
	storage *db.SQLCStore
	env     config.EnvVariables
}

func NewLoginService(storage *db.SQLCStore, env config.EnvVariables) *LoginService {
	return &LoginService{storage: storage, env: env}
}

type userDataModel struct {
	SocialId string `json:"social_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
}

func (service *LoginService) saveUserData(userData userDataModel, c context.Context) (db.User, error) {
	// does this user exist ?
	oldUser, getOldUserErr := service.storage.GetUserBySocialId(c, userData.SocialId)

	// if user doesn't exist, create it
	if getOldUserErr == sql.ErrNoRows {
		createUserParams := db.CreateUserParams{SocialID: userData.SocialId, Username: userData.Username, Email: userData.Email, Picture: userData.Picture, Provider: userData.Provider}
		userCreated, userCreatedErr := service.storage.CreateUser(c, createUserParams)
		if userCreatedErr != nil {
			err := fmt.Errorf("something went wrong creating a new user (%s provider) ----> %s", userData.Provider, userCreatedErr)
			return db.User{}, err
		}
		return userCreated, nil
	}

	// if user exists, update info (if neccesary)
	updateUserParams := db.UpdateUserParams{Username: userData.Username, Email: userData.Email, Pictue: userData.Picture, UserID: oldUser.UserID, OldUsername: oldUser.Username, OldEmail: oldUser.Email, OldPicture: oldUser.Picture}
	updatedUser, updatedUserErr := service.storage.UpdateUser(c, updateUserParams)
	if updatedUserErr == sql.ErrNoRows { // query didn't update the record and returned nothing
		return oldUser, nil
	} else if updatedUserErr != nil { // something happened executing the query
		err := fmt.Errorf("something went wrong updating new user (%s provider) ----> %s", userData.Provider, updatedUserErr)
		return db.User{}, err
	}

	return updatedUser, nil
}

func (service *LoginService) createSession(userId int64, userAgent string, c context.Context) (db.Session, error) {
	createSessionParams := db.GenerateSessionParams{UserIDSession: userId, Valid: true, UserAgent: userAgent}
	createdSession, sessionErr := service.storage.GenerateSession(c, createSessionParams)
	if sessionErr != nil {
		err := fmt.Errorf("couldn't create the session for user %d, error : %s", userId, sessionErr)
		return db.Session{}, err
	}
	return createdSession, nil
}
