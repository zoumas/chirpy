package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zoumas/chirpy/json/internal/database"
)

type AuthedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (app *App) WithAccessToken(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		authFields := strings.Fields(authHeader)
		if len(authFields) != 2 {
			respondWithError(w, http.StatusUnauthorized, "malformed Authorization header")
			return
		}

		authMethod := authFields[0]
		if authMethod != "Bearer" {
			respondWithError(
				w,
				http.StatusUnauthorized,
				fmt.Sprintf("Authorization method %q is not supported", authMethod),
			)
			return
		}

		tokenString := authFields[1]
		token, err := jwt.ParseWithClaims(
			tokenString,
			&jwt.RegisteredClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(app.Env.JwtSecret), nil
			},
		)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		issuer, err := token.Claims.GetIssuer()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if issuer == "chirpy-refresh" {
			respondWithError(w, http.StatusUnauthorized, "access token required")
			return
		}

		userIDString, err := token.Claims.GetSubject()
		if err != nil {
			respondWithError(
				w,
				http.StatusUnauthorized,
				fmt.Sprintf("failed to parse user ID from token: %s", err.Error()),
			)
			return
		}

		userID, err := strconv.Atoi(userIDString)
		if err != nil {
			respondWithError(
				w,
				http.StatusInternalServerError,
				fmt.Sprintf("failed to parse user ID: %s", err.Error()),
			)
			return
		}

		user, err := app.UserRepository.GetByID(userID)
		if err != nil {
			respondWithError(
				w,
				http.StatusUnauthorized,
				fmt.Sprintf("failed to retrieve user : %s", err.Error()),
			)
		}

		handler(w, r, user)
	}
}

type WithRefreshTokenParams struct {
	token  string
	userID int
}

func (app *App) WithRefreshToken(
	handler func(w http.ResponseWriter, r *http.Request, params WithRefreshTokenParams),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		authFields := strings.Fields(authHeader)
		if len(authFields) != 2 {
			respondWithError(w, http.StatusUnauthorized, "malformed Authorization header")
			return
		}

		authMethod := authFields[0]
		if authMethod != "Bearer" {
			respondWithError(
				w,
				http.StatusUnauthorized,
				fmt.Sprintf("Authorization method %q is not supported", authMethod),
			)
			return
		}

		tokenString := authFields[1]
		token, err := jwt.ParseWithClaims(
			tokenString,
			&jwt.RegisteredClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(app.Env.JwtSecret), nil
			},
		)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		issuer, err := token.Claims.GetIssuer()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if issuer == "chirpy-access" {
			respondWithError(w, http.StatusUnauthorized, "refresh token required")
			return
		}

		userIDString, err := token.Claims.GetSubject()
		if err != nil {
			respondWithError(
				w,
				http.StatusUnauthorized,
				fmt.Sprintf("failed to parse user ID from token: %s", err.Error()),
			)
			return
		}

		userID, err := strconv.Atoi(userIDString)
		if err != nil {
			respondWithError(
				w,
				http.StatusInternalServerError,
				fmt.Sprintf("failed to parse user ID: %s", err.Error()),
			)
			return
		}

		signedToken, err := token.SignedString([]byte(app.Env.JwtSecret))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		handler(w, r, WithRefreshTokenParams{token: signedToken, userID: userID})
	}
}

func (app *App) WithPolkaApiKey(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authFields := strings.Fields(authHeader)
		if len(authFields) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authMethod := authFields[0]
		if authMethod != "ApiKey" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		polkaApiKey := authFields[1]
		if polkaApiKey != app.Env.PolkaApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
