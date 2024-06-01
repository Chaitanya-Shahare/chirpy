package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
	RefreshToken     string `json:"refresh_token"`
}

// handlerLogin handles user login and JWT generation
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var u User

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

	user, err := cfg.DB.GetUserByEmail(u.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	userIDStr := strconv.Itoa(user.ID)

	tokenString, err := cfg.createJWT(userIDStr, 3600) // Access token expires in 1 hour
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT")
		return
	}

	refreshToken, err := cfg.createRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}

	// Store the refresh token with the user in the database
	err = cfg.DB.UpdateUserRefreshToken(user.ID, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error storing refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Email        string `json:"email"`
		ID           int    `json:"id"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		Email:        user.Email,
		ID:           user.ID,
		Token:        tokenString,
		RefreshToken: refreshToken,
	})
}

// createJWT creates a new JWT for the given user ID and expiration time
func (cfg *apiConfig) createJWT(userID string, expiresInSeconds int) (string, error) {
	expirationTime := time.Duration(expiresInSeconds) * time.Second

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expirationTime)),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign the token: %w", err)
	}

	return tokenString, nil
}

// createRefreshToken generates a new refresh token
func (cfg *apiConfig) createRefreshToken() (string, error) {
	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshToken := hex.EncodeToString(refreshTokenBytes)
	return refreshToken, nil
}
