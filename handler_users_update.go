package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(auth, "Bearer ")

	// Log the token string
	log.Printf("Token string: %s", tokenString)

	claims, err := cfg.validateJWT(tokenString)
	if err != nil {
		log.Printf("Token validation error: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Log the claims
	log.Printf("Claims: %+v", claims)

	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	// Log the subject before conversion
	log.Printf("Claims subject before conversion: %s", claims.Subject)

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Printf("Error converting subject to int: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update user (invalid ID)")
		return
	}

	// hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	_, err = cfg.DB.UpdateUser(id, u.Email, string(hashedPassword))

	if err != nil {
		// Handle error
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	}{
		Email: u.Email,
		ID:    id,
	})
}

func (cfg *apiConfig) validateJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	if !token.Valid {
		log.Println("Token is invalid")
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
