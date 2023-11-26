package app

import (
	"encoding/json"
	"net/http"

	"github.com/zoumas/chirpy/json/internal/database"
)

type JSONUserRepository struct {
	db *database.DB
}

func NewJSONUserRepository(db *database.DB) *JSONUserRepository {
	return &JSONUserRepository{
		db: db,
	}
}

func (r *JSONUserRepository) Create(email string) (database.User, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return database.User{}, err
	}

	// TODO: check if email is taken

	id := len(dbs.Users) + 1
	user := database.User{ID: id, Email: email}
	dbs.Users[id] = user

	err = r.db.Persist(dbs)
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email string `json:"email"`
	}
	body := RequestBody{}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := app.UserRepository.Create(body.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
