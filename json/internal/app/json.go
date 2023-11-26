package app

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, httpStatusCode int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, httpStatusCode int, errorMessage string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, httpStatusCode, ErrorResponse{errorMessage})
}