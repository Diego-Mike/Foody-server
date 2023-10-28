package mw

import (
	"context"
	"net/http"
	"strings"

	"github.com/Foody-App-Tech/Main-server/config"
	constants "github.com/Foody-App-Tech/Main-server/internal/constants"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
)

type GlobalMiddlewares struct {
	AccessTokenKey  string
	RefreshTokenKey string
	storage         *db.SQLCStore
}

func NewGlobalMiddlewareService(accessTokenKey, refreshTokenKey string, storage *db.SQLCStore) *GlobalMiddlewares {
	return &GlobalMiddlewares{storage: storage, AccessTokenKey: accessTokenKey, RefreshTokenKey: refreshTokenKey}
}

// FIXME: return message if user does not have access
func (service *GlobalMiddlewares) IdentifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get access token and refresh token
		getAccessToken := strings.Split(r.Header.Get("Authorization"), "Bearer")
		if len(getAccessToken) < 2 {
			config.ErrorResponse(w, "entidad no autorizada para hacer esta petición", nil, http.StatusUnauthorized)
			return
		}
		accessToken := strings.TrimSpace(getAccessToken[1])

		refreshTokenByCookie, refreshTokenErr := r.Cookie("refresh-token") // get refresh token
		if refreshTokenErr != nil {
			next.ServeHTTP(w, r)
			return
		}

		accessTokenDecoded, tokenErr := service.VerifyToken(accessToken, service.AccessTokenKey) // decode access token
		if tokenErr == ErrInvalidToken {
			config.ErrorResponse(w, "entidad no autorizada para hacer esta petición", nil, http.StatusUnauthorized)
			return
		}
		if accessTokenDecoded != nil {
			ctx := context.WithValue(r.Context(), constants.UserContextKey, accessTokenDecoded.UserData) // add user data to context (so we can access it)
			next.ServeHTTP(w, r.WithContext(ctx))                                                        // pass data to next middleware/handler
			return
		}

		// create new access token
		if (tokenErr == ErrExpiredToken) && (refreshTokenByCookie.Value != "") {
			newAccessToken, notAuthorized := service.ReIssueAccessToken(refreshTokenByCookie.Value, r.Context())

			if notAuthorized {
				config.ErrorResponse(w, "entidad no autorizada para hacer esta petición", nil, http.StatusUnauthorized)
				return
			}

			// FIXME: change to secure to true in prod
			// FIXME: how long should the cookie for access token and refresh token be
			// TODO: change cookies to right time in prod, leave this when developing
			http.SetCookie(w, &http.Cookie{Name: "AccessToken", Value: newAccessToken, Path: "/", Domain: "localhost", MaxAge: 900, HttpOnly: true, Secure: false}) // 15min

			decodedAccessToken, _ := service.VerifyToken(newAccessToken, service.AccessTokenKey)
			ctx := context.WithValue(r.Context(), constants.UserContextKey, decodedAccessToken.UserData)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})

}
