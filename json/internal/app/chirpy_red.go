package app

import (
	"encoding/json"
	"net/http"
)

func (app *App) SubscribeToChirpyRed(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	body := RequestBody{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		// This is really bad btw. If you can't decode a webhook's request body then
		// the webhook will be trying the same request you can't decode over and over again.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if body.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err = app.UserRepository.UpgradeToRed(body.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
