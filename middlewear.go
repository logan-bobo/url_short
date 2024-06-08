package main

import (
	"net/http"
	"strconv"
	"strings"
	"url-short/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

type authedHandeler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) authenticationMiddlewear(handler authedHandeler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			respondWithError(w, http.StatusBadRequest, "no auth header supplied")
			return
		}

		splitAuth := strings.Split(authHeader, " ")

		if len(splitAuth) == 0 {
			respondWithError(w, http.StatusBadRequest, "empty auth header")
		}

		if len(splitAuth) != 2 && splitAuth[0] != "Bearer" {
			respondWithError(w, http.StatusBadRequest, "invalid paremeters")
		}

		requestToken := splitAuth[1]

		claims := jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(
			requestToken,
			&claims,
			func(token *jwt.Token) (interface{}, error) { return []byte(apiCfg.JWTSecret), nil },
		)

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid jwt")
			return
		}

		issuer, err := token.Claims.GetIssuer()

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid jwt issuer")
			return
		}

		if issuer != "url-short-auth" {
			respondWithError(w, http.StatusUnauthorized, "invalid jwt issuer")
			return
		}

		userID, err := token.Claims.GetSubject()

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "could not get subject from jwt")
			return
		}

		userIDStr, err := strconv.Atoi(userID)

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid jwt subject")
		}

		user, err := apiCfg.DB.SelectUserByID(r.Context(), int32(userIDStr))

		handler(w, r, user)
	})
}
