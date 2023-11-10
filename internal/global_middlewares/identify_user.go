package mw

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/Foody-App-Tech/Main-server/config"
	constants "github.com/Foody-App-Tech/Main-server/internal/constants"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
)

type GlobalMiddlewares struct {
	AccessTokenKey   string
	AccessTokenTime  string
	RefreshTokenKey  string
	RefreshTokenTime string
	ApiKey           string
	SecureCookies    string
	storage          *db.SQLCStore
}

func NewGlobalMiddlewareService(accessTokenKey, accessTokenTime, refreshTokenKey, refreshTokenTime, apiKey, secureCookies string, storage *db.SQLCStore) *GlobalMiddlewares {
	return &GlobalMiddlewares{storage: storage, AccessTokenKey: accessTokenKey, AccessTokenTime: accessTokenTime, RefreshTokenKey: refreshTokenKey, RefreshTokenTime: refreshTokenTime, SecureCookies: secureCookies, ApiKey: apiKey}
}

func (service *GlobalMiddlewares) IdentifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get api key
		apiKey := r.Header.Get("foody-api-key")
		if (apiKey == "") || (apiKey != service.ApiKey) {
			config.ErrorResponse(w, "entidad no autorizada para hacer esta petición", nil, http.StatusUnauthorized)
			return
		}
		// log.Println("api key", apiKey)

		// get access token and refresh token
		accessTokenByCookie, accessTokenErr := r.Cookie("access-token")
		if accessTokenErr != nil {
			log.Println("no access token", accessTokenErr)
			config.ErrorResponse(w, "entidad no autorizada para hacer esta petición", nil, http.StatusUnauthorized)
			return
		}
		// log.Println("access token", accessTokenByCookie.Value)

		refreshTokenByCookie, refreshTokenErr := r.Cookie("refresh-token") // get refresh token
		if refreshTokenErr != nil {
			log.Println("no refresh token", refreshTokenErr)
			config.ErrorResponse(w, "entidad no autorizada para hacer esta petición", nil, http.StatusUnauthorized)
			return
		}

		accessTokenDecoded, tokenErr := service.VerifyToken(accessTokenByCookie.Value, service.AccessTokenKey) // decode access token
		if tokenErr == ErrInvalidToken {
			config.ErrorResponse(w, "entidad no autorizada para hacer esta petición", nil, http.StatusUnauthorized)
			return
		}
		// log.Println("we got here", accessTokenDecoded)
		if accessTokenDecoded != nil {
			ctx := context.WithValue(r.Context(), constants.UserContextKey, accessTokenDecoded.UserData) // add user data to context (so we can access it)
			next.ServeHTTP(w, r.WithContext(ctx))                                                        // pass data to next middleware/handler
			return
		}

		// create new access token
		if (strings.Contains(tokenErr.Error(), "token is expired")) && (refreshTokenByCookie.Value != "") {
			// TODO: update refresh token as well
			newAccessToken, notAuthorized := service.ReIssueAccessToken(w, refreshTokenByCookie.Value, r.Context())

			if notAuthorized {
				// log.Println("we got here")
				config.ErrorResponse(w, "entidad no autorizada para hacer esta petición", nil, http.StatusUnauthorized)
				return
			}

			decodedAccessToken, _ := service.VerifyToken(newAccessToken, service.AccessTokenKey)
			ctx := context.WithValue(r.Context(), constants.UserContextKey, decodedAccessToken.UserData)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// return error
		config.ErrorResponse(w, "entidad no autorizada para hacer esta petición", nil, http.StatusUnauthorized)
	})

}
