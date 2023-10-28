package mw

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUserData struct {
	Username string `json:"username"`
	Email    string `json:"Email"`
	SocialID string `json:"social_id"`
	UserID   int64  `json:"user_id"`
	Picture  string `json:"picture"`
}

// FIXME: include this stuff in the main struct ?¿
type JwtClaims struct {
	UserData JwtUserData
	jwt.RegisteredClaims
}

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// FIXME: test this
func (service *GlobalMiddlewares) CreateToken(userData JwtUserData, duration time.Duration, key string) (string, error) {

	tokenDuration := jwt.NewNumericDate(time.Now().Add(duration)) // how long will the token live

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{UserData: userData, RegisteredClaims: jwt.RegisteredClaims{
		Issuer:    "Foody users service",
		ExpiresAt: tokenDuration,
	}})

	signedToken, signedTokenErr := token.SignedString([]byte(key))
	if signedTokenErr != nil {
		err := fmt.Errorf("there was a problem signing token ----> %s", signedTokenErr)
		return "", err
	}

	return signedToken, nil

}

// this should be global
func (service *GlobalMiddlewares) CreateRefreshToken(user JwtUserData) (refreshToken string, err error) {

	// refresh token TODO: should this token expire at all ? should it last forever ?
	refreshTokenDuration := time.Hour * 2 // 1 month
	refreshToken, refreshTokenErr := service.CreateToken(user, refreshTokenDuration, service.RefreshTokenKey)
	if refreshTokenErr != nil {
		err = fmt.Errorf("there was a problem signing refresh token ----> %s", refreshTokenErr)
		return
	}

	return
}

// FIXME: test this
func (service *GlobalMiddlewares) VerifyToken(token, key string) (*JwtClaims, error) {

	// FIXME: check with key
	// check if jwt signing method is correct and return key
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(key), nil
	}

	// parse and validate the token
	jwtToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %s", err)
	}

	// check if token is valid
	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	// access claims
	claims, ok := jwtToken.Claims.(*JwtClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// check if token has expired
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrExpiredToken
	}

	return claims, nil

}

// FIXME: test this
func (service *GlobalMiddlewares) ReIssueAccessToken(refreshToken string, ctx context.Context) (string, bool) {
	// decode refresh token and check it
	refreshDecoded, refreshErr := service.VerifyToken(refreshToken, service.RefreshTokenKey)
	if refreshErr == ErrInvalidToken {
		return "", true
	}
	// create new refresh token ?¿
	if refreshErr == ErrExpiredToken {
		return "", true
	}

	// validate session
	session, sessionErr := service.storage.GetSession(ctx, refreshDecoded.UserData.UserID)
	if sessionErr != nil || !session.Valid {
		return "", true
	}

	// validate user
	user, userErr := service.storage.GetUserById(ctx, refreshDecoded.UserData.UserID)
	if userErr != nil {
		return "", true
	}

	// new access token
	createTokenParams := JwtUserData{Picture: user.Picture, Username: user.Username, Email: user.Email, SocialID: user.SocialID, UserID: user.UserID}
	accessToken, accessTokenErr := service.CreateToken(createTokenParams, time.Minute*15, service.AccessTokenKey)
	if accessTokenErr != nil {
		return "", true
	}

	return accessToken, false

}
