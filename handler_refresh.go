package main

import (
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	refreshtokenString := strings.TrimPrefix(auth, "Bearer ")

	user, err := cfg.DB.GetUserByRefreshToken(refreshtokenString)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	userIDStr := strconv.Itoa(user.ID)
	tokenString, err := cfg.createJWT(userIDStr, 3600)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	})

}
