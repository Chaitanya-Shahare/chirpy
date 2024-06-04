package main

import (
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpID)

	auth := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(auth, "Bearer ")

	claims, err := cfg.validateJWT(tokenString)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	userID, err := strconv.Atoi(claims.Subject)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't extract user ID")
		return
	}

	if chirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
