package login

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Foody-App-Tech/Main-server/config"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
)

type GoogleOauthService struct {
	storage *db.SQLCStore
	env     config.EnvVariables
}

func NewGoogleOauthService(storage *db.SQLCStore, env config.EnvVariables) *GoogleOauthService {
	return &GoogleOauthService{
		storage: storage,
		env:     env,
	}
}

type googleTokensModel struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IdToken      string `json:"id_token"`
}

func (service *GoogleOauthService) getGoogleUserTokens(w http.ResponseWriter, r *http.Request) (tokens googleTokensModel, err error) {
	// get google refresh and access token from user
	googleCode := r.URL.Query().Get("code")
	url := "https://oauth2.googleapis.com/token"

	qs := struct {
		code          string
		client_id     string
		client_secret string
		redirect_uri  string
		grant_type    string
	}{
		code:          googleCode,
		client_id:     service.env.GOOGLE_CLIENT_ID,
		client_secret: service.env.GOOGLE_CLIENT_SECRET,
		redirect_uri:  service.env.GOOGLE_REDIRECT_URI,
		grant_type:    service.env.GOOGLE_GRANT_TYPE,
	}
	values, qsErr := json.MarshalIndent(qs, "", "") // converting qs to []byte
	if qsErr != nil {
		err = fmt.Errorf("there was an error converting querystring values to bytes to get google user tokens ---> %s", qsErr)
		return
	}

	userTokensReq, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(values)) // creating the request
	if reqErr != nil {
		err = fmt.Errorf("there was an error making the request ---> %s", reqErr)
		return
	}
	userTokensReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	query := userTokensReq.URL.Query()
	query.Add("code", qs.code)
	query.Add("client_id", qs.client_id)
	query.Add("client_secret", qs.client_secret)
	query.Add("redirect_uri", qs.redirect_uri)
	query.Add("grant_type", qs.grant_type)
	userTokensReq.URL.RawQuery = query.Encode()

	userTokensClient := &http.Client{Timeout: time.Second * 10} // making the request
	googleTokensRes, googleTokensErr := userTokensClient.Do(userTokensReq)
	if googleTokensErr != nil {
		err = fmt.Errorf("there was an error making the request to get google user tokens ---> %s", googleTokensErr)
		return
	}

	readGoogleTokensRes, readGoogleTokensErr := io.ReadAll(googleTokensRes.Body) // transforming response to json
	if readGoogleTokensErr != nil {
		err = fmt.Errorf("there was an error reading google user tokens response ---> %s", readGoogleTokensErr)
		return
	}
	googleTokensJsonError := json.Unmarshal(readGoogleTokensRes, &tokens)
	if googleTokensJsonError != nil {
		err = fmt.Errorf("there was an error converting google user tokens to json ---> %s", googleTokensJsonError)
		return
	}

	if !strings.Contains(strconv.Itoa(googleTokensRes.StatusCode), "2") {
		log.Printf("google tokens response ----> %s", config.PrettyPrint(tokens))
		err = fmt.Errorf("status code from request to get google user tokens is not 200, something went wrong ----> %s", googleTokensRes.Status)
		return
	}

	return
}

type googleUserDataModel struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (service *GoogleOauthService) getGoogleUserData(userTokens googleTokensModel) (userData googleUserDataModel, err error) {
	googleUserDataURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", userTokens.AccessToken)

	googleUserReq, reqErr := http.NewRequest("GET", googleUserDataURL, nil)
	googleUserReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userTokens.IdToken))
	if reqErr != nil {
		err = fmt.Errorf("there was an error creating the request to get user data ----> %s", reqErr)
		return
	}

	googleUserDataClient := &http.Client{Timeout: time.Second * 10}
	googleUserDataRes, googleUserDataErr := googleUserDataClient.Do(googleUserReq)
	if googleUserDataErr != nil {
		err = fmt.Errorf("there was an error making the request to get user data ----> %s", reqErr)
		return
	}

	readGoogleUserDataRes, readGoogleUserDataErr := io.ReadAll(googleUserDataRes.Body)
	if readGoogleUserDataErr != nil {
		err = fmt.Errorf("something bad happened reading google user data response ----> %s", readGoogleUserDataErr)
		return
	}
	googleUserDataJsonErr := json.Unmarshal(readGoogleUserDataRes, &userData)
	if googleUserDataErr != nil {
		err = fmt.Errorf("something bad happened converting google user data response to json ----> %s", googleUserDataJsonErr)
		return
	}

	if !strings.Contains(strconv.Itoa(googleUserDataRes.StatusCode), "2") {
		log.Printf("google user data response ----> %s", config.PrettyPrint(userData))
		err = fmt.Errorf("status code from request to get google user data is not 200, something went wrong ----> %s", googleUserDataRes.Status)
		return
	}

	return
}

type userDataModel struct {
	SocialId string `json:"social_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
}

func (service *GoogleOauthService) saveUserData(userData userDataModel, c context.Context) (db.User, error) {
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

	// FIXME: bug, this updates all users
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

func (service *GoogleOauthService) createSession(userId int64, userAgent string, c context.Context) (db.Session, error) {
	createSessionParams := db.GenerateSessionParams{UserIDSession: userId, Valid: true, UserAgent: userAgent}
	createdSession, sessionErr := service.storage.GenerateSession(c, createSessionParams)
	if sessionErr != nil {
		err := fmt.Errorf("couldn't create the session for user %d, error : %s", userId, sessionErr)
		return db.Session{}, err
	}
	return createdSession, nil
}
