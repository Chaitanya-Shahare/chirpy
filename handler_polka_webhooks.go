package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type PolkaWebhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {

	// {
	// "event": "user.upgraded",
	// "data": {
	// "user_id": 3
	// }
	// }

	auth := r.Header.Get("Authorization")
	// tokenString := strings.TrimPrefix(auth, "Bearer ")
	apiKey := strings.TrimPrefix(auth, "ApiKey ")

	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API Key")
		return
	}

	var webhook PolkaWebhook

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&webhook); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	_, err := cfg.DB.GetUserByID(webhook.Data.UserID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	if webhook.Event == "user.upgraded" {
		err := cfg.DB.UpgradeUser(webhook.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user")
		}
		w.WriteHeader(http.StatusNoContent)
		// send a empty object
		w.Write([]byte("{}"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
