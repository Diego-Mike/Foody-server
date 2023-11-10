package login

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Foody-App-Tech/Main-server/config"
	mw "github.com/Foody-App-Tech/Main-server/internal/global_middlewares"
)

type LoginController struct {
	googleOauthService   *GoogleOauthService
	facebookOauthService *FacebookOauthService
	loginService         *LoginService
	globalHelpers        *mw.GlobalMiddlewares
}

func NewLoginController(googleOauthService *GoogleOauthService, facebookOauthService *FacebookOauthService, loginService *LoginService, globalHelpers *mw.GlobalMiddlewares) *LoginController {
	return &LoginController{googleOauthService: googleOauthService, facebookOauthService: facebookOauthService, loginService: loginService, globalHelpers: globalHelpers}
}

// FIXME: remove unnecesary logs when deploying
func (l *LoginController) googleLogin(w http.ResponseWriter, r *http.Request) {

	// get google user auth tokens
	// FIXME: change redirect url when something goes wrong
	googleUserTokens, err := l.googleOauthService.getGoogleUserTokens(w, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	// log.Printf("google user tokens %s", config.PrettyPrint(googleUserTokens))

	// get user with tokens
	googleUserData, err := l.googleOauthService.getGoogleUserData(googleUserTokens)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	// log.Printf("google user Data %s", config.PrettyPrint(googleUserData))

	// create-update user
	userData := userDataModel{SocialId: googleUserData.Id, Username: googleUserData.Name, Email: googleUserData.Email, Picture: googleUserData.Picture, Provider: "google"}
	user, err := l.loginService.saveUserData(userData, r.Context())
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	// log.Printf("google user %s", config.PrettyPrint(user))

	// create-update a session
	userAgent := r.UserAgent()
	_, err = l.loginService.createSession(user.UserID, userAgent, r.Context())
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	// log.Printf("google user session %s", config.PrettyPrint(session))

	// refresh token
	refreshTokenData := mw.JwtUserData{UserID: user.UserID}
	_, err = l.globalHelpers.CreateRefreshOrAccessToken(w, refreshTokenData, l.googleOauthService.env.REFRESH_TOKEN_TIME, l.googleOauthService.env.REFRESH_TOKEN_KEY, "refresh-token")
	if err != nil {
		log.Println("problem creating refresh token when logging in", err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	// log.Println("refresh token", refreshToken)

	// access token
	accessTokenData := mw.JwtUserData{Username: user.Username, Email: user.Email, SocialID: user.SocialID, UserID: user.UserID, Picture: user.Picture}
	_, err = l.globalHelpers.CreateRefreshOrAccessToken(w, accessTokenData, l.googleOauthService.env.ACCESS_TOKEN_TIME, l.googleOauthService.env.ACCESS_TOKEN_KEY, "access-token")
	if err != nil {
		log.Println("problem creating refresh token when logging in", err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	// log.Println("access token", accessToken)

	// redirect back to client
	http.Redirect(w, r, "http://localhost:3000", http.StatusSeeOther)
}

func (l *LoginController) facebookLogin(w http.ResponseWriter, r *http.Request) {

	// get access token
	facebookToken, err := l.facebookOauthService.getFacebookAccessToken(r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/facebook_login/error", http.StatusSeeOther)
		return
	}
	// log.Printf("facebook access token %s", config.PrettyPrint(facebookToken))

	// get user data
	facebookUser, err := l.facebookOauthService.getFacebookUserData(facebookToken)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/facebook_login/error", http.StatusSeeOther)
	}
	// log.Printf("facebook user data %s", config.PrettyPrint(facebookUser))

	// create-update-user
	userData := userDataModel{SocialId: facebookUser.Id, Username: facebookUser.Name, Email: facebookUser.Email, Picture: facebookUser.Picture.Data.Url, Provider: "facebook"}
	user, err := l.loginService.saveUserData(userData, r.Context())
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/facebook_login/error", http.StatusSeeOther)
		return
	}
	// log.Printf("facebook user %s", config.PrettyPrint(user))

	// create-update session
	userAgent := r.UserAgent()
	_, err = l.loginService.createSession(user.UserID, userAgent, r.Context())
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/facebook_login/error", http.StatusSeeOther)
		return
	}
	// log.Printf("facebook user session %s", config.PrettyPrint(session))

	// refresh token
	refreshTokenData := mw.JwtUserData{UserID: user.UserID}
	_, err = l.globalHelpers.CreateRefreshOrAccessToken(w, refreshTokenData, l.facebookOauthService.env.REFRESH_TOKEN_TIME, l.facebookOauthService.env.REFRESH_TOKEN_KEY, "refresh-token")
	if err != nil {
		log.Println("problem creating refresh token when logging in", err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	// log.Println("refresh token", refreshToken)

	// access token
	accessTokenData := mw.JwtUserData{Username: user.Username, Email: user.Email, SocialID: user.SocialID, UserID: user.UserID, Picture: user.Picture}
	_, err = l.globalHelpers.CreateRefreshOrAccessToken(w, accessTokenData, l.facebookOauthService.env.ACCESS_TOKEN_TIME, l.facebookOauthService.env.ACCESS_TOKEN_KEY, "access-token")
	if err != nil {
		log.Println("problem creating refresh token when logging in", err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	// log.Println("access token", accessToken)

	// redirect back to client
	http.Redirect(w, r, "http://localhost:3000", http.StatusSeeOther)

}

func (l *LoginController) accessToken(w http.ResponseWriter, r *http.Request) {

	// check api key
	apiKey := r.Header.Get("foody-api-key")
	// log.Println("api key", apiKey)
	if (apiKey == "") || (apiKey != l.globalHelpers.ApiKey) {
		// log.Println("error with api key")
		config.ErrorResponse(w, "Entidad no autorizada para realizar esta petición !", nil, http.StatusUnauthorized)
		return
	}

	// get refresh token
	refreshToken, err := r.Cookie("refresh-token")
	// log.Println("refresh token", refreshToken)
	if err != nil {
		config.ErrorResponse(w, "Usuario no posee credenciales, inicia sesión !", nil, http.StatusUnauthorized)
		return
	}

	// create new access token
	accessToken, notAuthorized := l.globalHelpers.ReIssueAccessToken(w, refreshToken.Value, r.Context())
	if notAuthorized {
		config.ErrorResponse(w, "No se le puedo dar acceso a este usuario !", nil, http.StatusUnauthorized)
		return
	}

	expiryTime, _ := strconv.Atoi(l.globalHelpers.AccessTokenTime)
	secureCookies, _ := strconv.ParseBool(l.globalHelpers.SecureCookies)
	cookie := &http.Cookie{Name: "access-token", Value: accessToken, Path: "/", Domain: "localhost", MaxAge: expiryTime, HttpOnly: true, Secure: secureCookies} // for prod 15min - for dev 5min
	http.SetCookie(w, cookie)

	// send response
	rspData := struct {
		AccessToken string `json:"access_token"`
		Msg         string `json:"msg"`
	}{
		AccessToken: accessToken,
		Msg:         "Token generado exitosamente!",
	}
	response := config.ClientResponse{Rsp: rspData}
	config.WriteResponse(w, http.StatusOK, response)
}

func (l *LoginController) logout(w http.ResponseWriter, r *http.Request) {

	secureCookies, _ := strconv.ParseBool(l.globalHelpers.SecureCookies)

	//logout
	refreshTokenCookie := &http.Cookie{Name: "refresh-token", Value: "", MaxAge: -1, Path: "/", Domain: "localhost", Secure: secureCookies, HttpOnly: true}
	http.SetCookie(w, refreshTokenCookie)
	accessTokenCookie := &http.Cookie{Name: "access-token", Value: "", MaxAge: -1, Path: "/", Domain: "localhost", Secure: secureCookies, HttpOnly: true}
	http.SetCookie(w, accessTokenCookie)

	// send response
	response := config.ClientResponse{Rsp: "Successfull logout"}
	config.WriteResponse(w, http.StatusOK, response)

}
