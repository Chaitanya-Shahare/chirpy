package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var u User = User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if u.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if u.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	user, err := cfg.DB.CreateUser(u.Email, string(hashedPassword))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		Email       string `json:"email"`
		ID          int    `json:"id"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		Email:       user.Email,
		ID:          user.ID,
		IsChirpyRed: user.IsChirpyRed,
	})

}
