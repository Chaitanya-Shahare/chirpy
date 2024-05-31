package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

type apiConfig struct {
	fileserverHits int
	mu             sync.Mutex
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.mu.Lock()
		cfg.fileserverHits++
		cfg.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.mu.Lock()
	hits := cfg.fileserverHits
	cfg.mu.Unlock()
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, " <html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html> ", hits)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.mu.Lock()
	cfg.fileserverHits = 0
	cfg.mu.Unlock()
	fmt.Fprint(w, "Hits reset to 0")
}

type errorReturnVals struct {
	// the key will be the name of struct field unless you give it an explicit JSON tag
	Error string `json:"error"`
}

func somethingWentWrong(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(500)
	w.Header().Set("Content-Type", "application/json")

	errBody := errorReturnVals{
		Error: "Something went wrong",
	}

	if dat, err := json.Marshal(errBody); err == nil {

		w.Write(dat)
		return
	}

	w.Write([]byte("error"))
	return
}

func addChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)

		somethingWentWrong(w, r)

		return

	}

	if len(params.Body) > 140 {
		errBody := errorReturnVals{
			Error: "Chirp is too long",
		}

		if dat, err := json.Marshal(errBody); err != nil {
			somethingWentWrong(w, r)
			return
		} else {

			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			w.Write(dat)
		}

		return
	}

	s := strings.Split(params.Body, " ")

	for i, wo := range s {
		word := strings.ToLower(wo)
		if word == "kerfuffle" || word == "sharbert" || word == "fornax" {
			// isProfane = true
			s[i] = "****"
		}
	}

	type profanedReturnVals struct {
		CleanedBody string `json:"body"`
	}

	respBody := profanedReturnVals{
		CleanedBody: strings.Join(s, " "),
	}

	dat, err := json.Marshal(respBody)

	if err != nil {
		somethingWentWrong(w, r)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
	return

}

func main() {

	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	fs := http.FileServer(http.Dir("."))
	handler := apiCfg.middlewareMetricsInc(fs)
	mux.Handle("/app/", http.StripPrefix("/app", handler))

	mux.HandleFunc("GET /api/healthz", healthzHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)

	mux.HandleFunc("POST /api/chirps", addChirpsHandler)

	// mux.HandleFunc("GET /api/chirps", getChirpsHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server at port 8080")

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
