package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/Chaitanya-Shahare/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	// Get the author_id query parameter from the request
	authorID := r.URL.Query().Get("author_id")

	// Check if the sort query parameter is provided
	sortParam := r.URL.Query().Get("sort")
	if sortParam != "asc" && sortParam != "desc" {
		sortParam = "asc" // Default sort order is ascending
	}

	// Retrieve chirps based on author_id if provided
	var dbChirps []database.Chirp
	var err error
	if authorID != "" {
		// Convert authorID to integer
		var authorIDInt int
		authorIDInt, err = strconv.Atoi(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		// Retrieve chirps by author ID
		dbChirps, err = cfg.DB.GetChirpsByAuthorID(authorIDInt)
	} else {
		// Retrieve all chirps
		dbChirps, err = cfg.DB.GetChirps()
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	// Convert retrieved chirps into the desired format
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	// Sort chirps based on sortParam
	if sortParam == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}

	// Respond with JSON
	respondWithJSON(w, http.StatusOK, chirps)
}
