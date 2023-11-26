package app

import (
	"cmp"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/zoumas/chirpy/json/internal/database"
)

type JSONChirpRepository struct {
	db *database.DB
}

func NewJSONChirpResository(db *database.DB) *JSONChirpRepository {
	return &JSONChirpRepository{
		db: db,
	}
}

// Create creates a new Chirp from a given body and stores it in the database; auto-incrementing the ID.
func (r *JSONChirpRepository) Create(body string) (database.Chirp, error) {
	dbs, err := r.db.Load()
	if err != nil {
		return database.Chirp{}, err
	}

	id := len(dbs.Chirps) + 1
	chirp := database.Chirp{ID: id, Body: body}
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

type ChirpErr string

func (e ChirpErr) Error() string {
	return string(e)
}

const (
	ErrTooLong = ChirpErr("Chirp is too long")
	ErrEmpty   = ChirpErr("Chirp is empty")
)

func ValidateChirpLength(body string) error {
	const MaxChirpLength = 140

	switch l := len(body); {
	case l == 0:
		return ErrEmpty
	case l > MaxChirpLength:
		return ErrTooLong
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

func (app *App) CreateChirp(w http.ResponseWriter, r *http.Request) {
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

	chirp, err := app.ChirpRepository.Create(cleanedBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (app *App) GetChirps(w http.ResponseWriter, r *http.Request) {
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
