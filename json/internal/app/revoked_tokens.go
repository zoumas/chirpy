package app

import (
	"net/http"

	"github.com/zoumas/chirpy/json/internal/database"
)

type JSONRevokedTokensRepository struct {
	db *database.DB
}

func NewJSONRevokedTokensRepository(db *database.DB) *JSONRevokedTokensRepository {
	return &JSONRevokedTokensRepository{db: db}
}

func (r *JSONRevokedTokensRepository) Revoke(token string) error {
	dbs, err := r.db.Load()
	if err != nil {
		return err
	}

	dbs.RevokedTokens[token] = struct{}{}

	err = r.db.Persist(dbs)
	if err != nil {
		return err
	}

	return nil
}

func (r *JSONRevokedTokensRepository) IsRevoked(token string) (bool, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return false, err
	}

	_, ok := dbs.RevokedTokens[token]
	return ok, nil
}

func (app *App) Revoke(w http.ResponseWriter, r *http.Request, params WithRefreshTokenParams) {
	err := app.RevokedTokensRepository.Revoke(params.token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) Refresh(w http.ResponseWriter, r *http.Request, params WithRefreshTokenParams) {
	ok, err := app.RevokedTokensRepository.IsRevoked(params.token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if ok {
		respondWithError(w, http.StatusUnauthorized, "token is revoked")
		return
	}

	accessToken := NewAccessToken(params.userID)
	signedAccessToken, err := accessToken.SignedString([]byte(app.Env.JwtSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type ResponseBody struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, ResponseBody{signedAccessToken})
}
