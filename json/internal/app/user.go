package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zoumas/chirpy/json/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type UserErr string

func (e UserErr) Error() string {
	return string(e)
}

const (
	ErrUserNotFound   = UserErr("User not found")
	ErrUserEmailTaken = UserErr("Email taken")
)

type JSONUserRepository struct {
	db *database.DB
}

func NewJSONUserRepository(db *database.DB) *JSONUserRepository {
	return &JSONUserRepository{db: db}
}

func (r *JSONUserRepository) Create(params database.CreateUserParams) (database.User, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return database.User{}, err
	}

	// TODO: check if email is taken
	if _, err := r.GetByEmail(params.Email); err == nil {
		return database.User{}, ErrUserEmailTaken
	}

	id := len(dbs.Users) + 1
	user := database.User{
		ID:       id,
		Email:    params.Email,
		Password: params.Password,
	}
	dbs.Users[id] = user

	err = r.db.Persist(dbs)
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}

func (r *JSONUserRepository) GetByEmail(email string) (database.User, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return database.User{}, err
	}

	for _, user := range dbs.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return database.User{}, ErrUserNotFound
}

func (r *JSONUserRepository) GetByID(id int) (database.User, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return database.User{}, err
	}

	user, ok := dbs.Users[id]
	if !ok {
		return database.User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *JSONUserRepository) Update(
	id int,
	params database.UpdateUserParams,
) (database.User, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return database.User{}, err
	}

	user, ok := dbs.Users[id]
	if !ok {
		return database.User{}, ErrUserNotFound
	}

	user.Email = params.Email
	user.Password = params.Password

	dbs.Users[user.ID] = user

	err = r.db.Persist(dbs)
	if err != nil {
		return database.User{}, err
	}

	return user, nil
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	body := RequestBody{}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := app.UserRepository.Create(database.CreateUserParams{
		Email:    body.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type ResponseBody struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	respondWithJSON(w, http.StatusCreated, ResponseBody{ID: user.ID, Email: user.Email})
}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	body := RequestBody{}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := app.UserRepository.GetByEmail(body.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	const twentyFourHrs = 24 * 60 * 60
	var expiresInSeconds int = twentyFourHrs
	if body.ExpiresInSeconds != 0 && body.ExpiresInSeconds < twentyFourHrs {
		expiresInSeconds = body.ExpiresInSeconds
	}

	expiresIn, err := time.ParseDuration(fmt.Sprintf("%d", expiresInSeconds) + "s")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create JWT")
		return
	}
	token := CreateJWT(user.ID, expiresIn)

	signedToken, err := token.SignedString([]byte(app.Env.JwtSecret))
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("failed to create JWT: %s", err),
		)
		return
	}

	type ResponseBody struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}
	respondWithJSON(
		w,
		http.StatusOK,
		ResponseBody{ID: user.ID, Email: user.Email, Token: signedToken},
	)
}

func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request, user database.User) {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	body := RequestBody{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if body.Email != user.Email {
		_, err := app.UserRepository.GetByEmail(body.Email)
		if err == ErrUserEmailTaken {
			respondWithError(w, http.StatusBadRequest, ErrUserEmailTaken.Error())
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedUser, err := app.UserRepository.Update(user.ID, database.UpdateUserParams{
		Email:    body.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type ResponseBody struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	respondWithJSON(w, http.StatusOK, ResponseBody{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
	})
}
