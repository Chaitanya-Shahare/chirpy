package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	auth := r.Header.Get("Authorization")
	refreshtokenString := strings.TrimPrefix(auth, "Bearer ")

	user, err := cfg.DB.GetUserByRefreshToken(refreshtokenString)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	err = cfg.DB.DeleteRefreshToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
