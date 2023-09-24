package login

import (
	"log"
	"net/http"

	"github.com/Foody-App-Tech/Main-server/config"
	mw "github.com/Foody-App-Tech/Main-server/internal/global_middlewares"
)

type LoginController struct {
	googleOauthService   *GoogleOauthService
	facebookOauthService *FacebookOauthService
	globalHelpers        *mw.GlobalMiddlewares
}

func NewLoginController(googleOauthService *GoogleOauthService, facebookOauthService *FacebookOauthService, globalHelpers *mw.GlobalMiddlewares) *LoginController {
	return &LoginController{googleOauthService: googleOauthService, facebookOauthService: facebookOauthService, globalHelpers: globalHelpers}
}

// FIXME: remove unnecesary logs when deploying
func (l *LoginController) googleLogin(w http.ResponseWriter, r *http.Request) {

	// FIXME: bug, our api returns wrong user data

	// get google user auth tokens
	googleUserTokens, err := l.googleOauthService.getGoogleUserTokens(w, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	log.Printf("google user tokens %s", config.PrettyPrint(googleUserTokens))

	// get user with tokens
	googleUserData, err := l.googleOauthService.getGoogleUserData(googleUserTokens)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	log.Printf("google user Data %s", config.PrettyPrint(googleUserData))

	// create-update user
	userData := userDataModel{SocialId: googleUserData.Id, Username: googleUserData.Name, Email: googleUserData.Email, Picture: googleUserData.Picture, Provider: "google"}
	user, err := l.googleOauthService.saveUserData(userData, r.Context())
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	log.Printf("google user %s", config.PrettyPrint(user))

	// create-update a session
	userAgent := r.UserAgent()
	session, err := l.googleOauthService.createSession(user.UserID, userAgent, r.Context())
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	log.Printf("google user session %s", config.PrettyPrint(session))

	// refresh token
	refreshToken, err := l.globalHelpers.CreateRefreshToken(mw.JwtUserData{Username: user.Username, Email: user.Email, SocialID: user.SocialID, UserId: user.UserID, Picture: user.Picture})
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	log.Printf("refreshToken : %s", refreshToken)

	// set cookies
	// FIXME: change to secure to true in prod
	// FIXME: how long should the cookie for access token and refresh token be
	// TODO: change cookies to right time in prod, leave this when developing
	http.SetCookie(w, &http.Cookie{Name: "refresh-token", Value: refreshToken, Path: "/", Domain: "localhost", MaxAge: 7200, HttpOnly: true, Secure: false}) // for prod 1month - for dev 2h

	// redirect back to client
	http.Redirect(w, r, "http://localhost:3000", http.StatusSeeOther)
}

func (l *LoginController) facebookLogin(w http.ResponseWriter, r *http.Request) {

	// FIXME: organize folder
	// FIXME: re-utilize functions, like, create a function to make a request to certain service

	// get access token
	facebookToken, err := l.facebookOauthService.getFacebookAccessToken(r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/facebook_login/error", http.StatusSeeOther)
		return
	}
	log.Printf("facebook access token %s", config.PrettyPrint(facebookToken))

	// get user data
	facebookUser, err := l.facebookOauthService.getFacebookUserData(facebookToken)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/facebook_login/error", http.StatusSeeOther)
	}
	log.Printf("facebook user data %s", config.PrettyPrint(facebookUser))

	// create-update-user
	userData := userDataModel{SocialId: facebookUser.Id, Username: facebookUser.Name, Email: facebookUser.Email, Picture: facebookUser.Picture.Data.Url, Provider: "facebook"}
	user, err := l.googleOauthService.saveUserData(userData, r.Context())
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/facebook_login/error", http.StatusSeeOther)
		return
	}
	log.Printf("facebook user %s", config.PrettyPrint(user))

	// create-update session
	userAgent := r.UserAgent()
	session, err := l.googleOauthService.createSession(user.UserID, userAgent, r.Context())
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/facebook_login/error", http.StatusSeeOther)
		return
	}
	log.Printf("facebook user session %s", config.PrettyPrint(session))

	// refresh token
	refreshToken, err := l.globalHelpers.CreateRefreshToken(mw.JwtUserData{Username: user.Username, Email: user.Email, SocialID: user.SocialID, UserId: user.UserID, Picture: user.Picture})
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "http://localhost:3000/google_login/error", http.StatusSeeOther)
		return
	}
	log.Printf("refreshToken : %s", refreshToken)

	// set cookies
	// FIXME: change to secure to true in prod
	// FIXME: how long should the cookie for access token and refresh token be
	// TODO: change cookies to right time in prod, leave this when developing
	http.SetCookie(w, &http.Cookie{Name: "refresh-token", Value: refreshToken, Path: "/", Domain: "localhost", MaxAge: 7200, HttpOnly: true, Secure: false}) // for prod 1month - for dev 2h

	// redirect back to client
	http.Redirect(w, r, "http://localhost:3000", http.StatusSeeOther)

}

func (l *LoginController) accessToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token
	refreshToken, err := r.Cookie("refresh-token")
	if err != nil {
		config.ErrorResponse(w, "user does not have credentials, log in", http.StatusUnauthorized)
		return
	}

	// create new access token
	accessToken := l.globalHelpers.ReIssueAccessToken(refreshToken.Value, r.Context())
	if accessToken == "" {
		config.ErrorResponse(w, "can't create authorization token for this user", http.StatusUnauthorized)
		return
	}

	// send access token
	response := config.ClientResponse{Error: false, Message: "access token generated successfully", Data: accessToken}
	config.WriteResponse(w, http.StatusOK, response)
}

func (l *LoginController) logout(w http.ResponseWriter, r *http.Request) {

	// TODO: change secure to true when deploying to prod
	//logout
	cookie := &http.Cookie{Name: "refresh-token", Value: "", MaxAge: -1, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true}
	http.SetCookie(w, cookie)

	// successfull response
	response := config.ClientResponse{Message: "Successfull logout"}
	config.WriteResponse(w, http.StatusOK, response)

}
