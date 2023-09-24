package login

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Foody-App-Tech/Main-server/config"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
)

type FacebookOauthService struct {
	storage *db.SQLCStore
	env     config.EnvVariables
}

func NewFacebookOauthService(storage *db.SQLCStore, env config.EnvVariables) *FacebookOauthService {
	return &FacebookOauthService{
		storage: storage,
		env:     env,
	}
}

type facebookTokenModel struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int32  `json:"expires_in"`
}

func (service *FacebookOauthService) getFacebookAccessToken(r *http.Request) (token facebookTokenModel, err error) {

	facebookCode := r.URL.Query().Get("code")
	url := "https://graph.facebook.com/v17.0/oauth/access_token"
	qs := struct {
		client_id     string
		redirect_uri  string
		client_secret string
		code          string
	}{
		client_id:     service.env.FACEBOOK_CLIENT_ID,
		redirect_uri:  service.env.FACEBOOK_REDIRECT_URI,
		client_secret: service.env.FACEBOOK_CLIENT_SECRET,
		code:          facebookCode,
	}

	facebookAccessTokenReq, err := http.NewRequest("GET", url, nil) // creating request
	if err != nil {
		err = fmt.Errorf("there was an error creating the request to get facebook user tokens ---> %s", err)
		return
	}
	query := facebookAccessTokenReq.URL.Query()
	query.Add("client_id", qs.client_id)
	query.Add("redirect_uri", qs.redirect_uri)
	query.Add("client_secret", qs.client_secret)
	query.Add("code", qs.code)
	facebookAccessTokenReq.URL.RawQuery = query.Encode()

	facebookUserTokenClient := &http.Client{Timeout: time.Second * 10} // making request
	facebookUserTokenRes, err := facebookUserTokenClient.Do(facebookAccessTokenReq)
	if err != nil {
		err = fmt.Errorf("there was an error making the request to get facebook user tokens ---> %s", err)
		return
	}

	facebookTokenRes, err := io.ReadAll(facebookUserTokenRes.Body) // transforming response to json
	if err != nil {
		err = fmt.Errorf("there was an error transforming facebook user tokens to json ---> %s", err)
		return
	}
	err = json.Unmarshal(facebookTokenRes, &token)
	if err != nil {
		err = fmt.Errorf("there was an error converting facebook user tokens to json ---> %s", err)
		return
	}

	return
}

type facebookUserDataModel struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture struct {
		Data struct {
			Height       int16  `json:"height"`
			IsSilhouette bool   `json:"is_silhouette"`
			Url          string `json:"url"`
			Width        int16  `json:"width"`
		}
	}
	Email string `json:"email"`
}

func (service *FacebookOauthService) getFacebookUserData(facebookToken facebookTokenModel) (facebookUser facebookUserDataModel, err error) {

	url := "https://graph.facebook.com/me"
	qs := struct {
		fields       string
		access_token string
	}{
		fields:       "id,name,email,picture",
		access_token: facebookToken.AccessToken,
	}

	facebookUserReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("there was an error making the request to get facebook user data ---> %s", err)
		return
	}
	query := facebookUserReq.URL.Query()
	query.Add("fields", qs.fields)
	query.Add("access_token", qs.access_token)
	facebookUserReq.URL.RawQuery = query.Encode()

	facebookUserDataClient := &http.Client{Timeout: time.Second * 10}
	facebookUserDataRes, err := facebookUserDataClient.Do(facebookUserReq)
	if err != nil {
		err = fmt.Errorf("there was an error making the request to get facebook user data ---> %s", err)
		return
	}

	facebookUserDataToJson, err := io.ReadAll(facebookUserDataRes.Body) // transforming response to json
	if err != nil {
		err = fmt.Errorf("there was an error transforming facebook user data to json ---> %s", err)
		return
	}
	err = json.Unmarshal(facebookUserDataToJson, &facebookUser)
	if err != nil {
		err = fmt.Errorf("there was an error transforming facebook user data to json ---> %s", err)
		return
	}

	return
}
