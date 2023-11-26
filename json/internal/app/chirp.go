package app

import (
	"cmp"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/zoumas/chirpy/json/internal/database"
)

type ChirpErr string

func (e ChirpErr) Error() string {
	return string(e)
}

const (
	ErrChirpTooLong   = ChirpErr("Chirp is too long")
	ErrChirpEmpty     = ChirpErr("Chirp is empty")
	ErrChirpNotFound  = ChirpErr("Chirp not found")
	ErrChirpNotAuthor = ChirpErr("Chirp is not owned by this user")
)

type JSONChirpRepository struct {
	db *database.DB
}

func NewJSONChirpResository(db *database.DB) *JSONChirpRepository {
	return &JSONChirpRepository{db: db}
}

// Create creates a new Chirp from a given body and stores it in the database; auto-incrementing the ID.
func (r *JSONChirpRepository) Create(params database.CreateChirpParams) (database.Chirp, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return database.Chirp{}, err
	}

	id := len(dbs.Chirps) + 1
	chirp := database.Chirp{ID: id, Body: params.Body, UserID: params.UserID}
	dbs.Chirps[id] = chirp

	err = r.db.Persist(dbs)
	if err != nil {
		return database.Chirp{}, err
	}

	return chirp, nil
}

// GetAll retrieves all the chirps from the database
func (r *JSONChirpRepository) GetAll() ([]database.Chirp, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return nil, err
	}

	chirps := make([]database.Chirp, 0, len(dbs.Chirps))
	for _, chirp := range dbs.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (r *JSONChirpRepository) GetByID(id int) (database.Chirp, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return database.Chirp{}, err
	}

	chirp, ok := dbs.Chirps[id]
	if !ok {
		return database.Chirp{}, ErrChirpNotFound
	}
	return chirp, nil
}

func (r *JSONChirpRepository) Delete(params database.DeleteChirpParams) error {
	dbs, err := r.db.Load()
	if err != nil {
		return err
	}

	chirp, ok := dbs.Chirps[params.ID]
	if !ok {
		return ErrChirpNotFound
	}

	if chirp.UserID != params.UserID {
		return ErrChirpNotAuthor
	}

	delete(dbs.Chirps, params.ID)
	return r.db.Persist(dbs)
}

func ValidateChirpLength(body string) error {
	const MaxChirpLength = 140

	switch l := len(body); {
	case l == 0:
		return ErrChirpEmpty
	case l > MaxChirpLength:
		return ErrChirpTooLong
	}

	return nil
}

func CleanChirpBody(
	body string,
	profane map[string]struct{},
	replaceWith string,
) (cleanedBody string) {
	split := strings.Split(body, " ")

	for i, word := range split {
		word = strings.ToLower(word)
		if _, ok := profane[word]; ok {
			split[i] = replaceWith
		}
	}

	return strings.Join(split, " ")
}

func (app *App) CreateChirp(w http.ResponseWriter, r *http.Request, user database.User) {
	type RequestBody struct {
		Body string `json:"body"`
	}
	body := RequestBody{}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ValidateChirpLength(body.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	token := struct{}{}
	profane := map[string]struct{}{
		"kerfuffle": token,
		"sharbert":  token,
		"fornax":    token,
	}
	replaceWith := "****"
	cleanedBody := CleanChirpBody(body.Body, profane, replaceWith)

	chirp, err := app.ChirpRepository.Create(database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (app *App) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := app.ChirpRepository.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve chirps")
		return
	}

	slices.SortStableFunc(chirps, func(a, b database.Chirp) int {
		return cmp.Compare(a.ID, b.ID)
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (app *App) GetChirpByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "missing url parameter")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse url parameter")
		return
	}

	chirp, err := app.ChirpRepository.GetByID(id)
	if err != nil {
		if err == ErrChirpNotFound {
			respondWithError(w, http.StatusNotFound, ErrChirpEmpty.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (app *App) DeleteChirp(w http.ResponseWriter, r *http.Request, user database.User) {
	idString := chi.URLParam(r, "id")
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "missing url parameter")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse url parameter")
		return
	}

	err = app.ChirpRepository.Delete(database.DeleteChirpParams{ID: id, UserID: user.ID})
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
