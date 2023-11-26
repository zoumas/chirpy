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

func (app *App) WithJWT(handler AuthedHandler) http.HandlerFunc {
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
