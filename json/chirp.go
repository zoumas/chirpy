package main

import (
	"encoding/json"
	"net/http"
)

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

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
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

	type ValidResponse struct {
		Valid bool `json:"valid"`
	}
	respondWithJSON(w, http.StatusOK, ValidResponse{true})
}
