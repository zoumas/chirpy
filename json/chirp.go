package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

	token := struct{}{}
	profane := map[string]struct{}{
		"kerfuffle": token,
		"sharbert":  token,
		"fornax":    token,
	}
	replaceWith := "****"
	cleanedBody := CleanChirpBody(body.Body, profane, replaceWith)

	type CleanedBodyResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respondWithJSON(w, http.StatusOK, CleanedBodyResponse{cleanedBody})
}
