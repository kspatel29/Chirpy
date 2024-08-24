package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kspatel29/chirpy/internal/database"
)

var (
        profaneWords = []string{
                "kerfuffle",
                "sharbert",
                "fornax",
        }
)

func main() {
        db, err := database.NewDB("database.json")
        if err != nil {
                panic(err)
        }

        mux := http.NewServeMux()

        mux.HandleFunc("/api/chirps", func(w http.ResponseWriter, r *http.Request) {
                switch r.Method {
                case http.MethodPost:
                        handleCreateChirp(w, r, db)
                case http.MethodGet:
                        handleGetChirps(w, r, db)
                default:
                        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
                }
        })

        server := &http.Server{
                Addr:    ":8080",
                Handler: mux,
        }

        if err := server.ListenAndServe(); err != nil {
                panic(err)
        }
}

func handleCreateChirp(w http.ResponseWriter, r *http.Request, db *database.DB) {
        var chirpRequest struct {
                Body string `json:"body"`
        }

        if err := json.NewDecoder(r.Body).Decode(&chirpRequest); err != nil {
                respondWithError(w, http.StatusBadRequest, "Invalid request payload")
                return
        }

        // Validate the chirp length
        if len(chirpRequest.Body) > 140 {
                respondWithError(w, http.StatusBadRequest, "Chirp is too long")
                return
        }

        // Clean the chirp from profane words
        cleanedBody := cleanProfanity(chirpRequest.Body)

        // Create the chirp and save it to the database
        chirp, err := db.CreateChirp(cleanedBody)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Could not save chirp")
                return
        }

        // Respond with the created chirp
        respondWithJSON(w, http.StatusCreated, chirp)
}

func handleGetChirps(w http.ResponseWriter, r *http.Request, db *database.DB) {
        chirps, err := db.GetChirps()
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps")
                return
        }

        respondWithJSON(w, http.StatusOK, chirps)
}

func cleanProfanity(text string) string {
        words := strings.Split(text, " ")
        for i, word := range words {
                for _, profaneWord := range profaneWords {
                        if strings.ToLower(word) == profaneWord {
                                words[i] = "****"
                        }
                }
        }
        return strings.Join(words, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
        respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
        response, _ := json.Marshal(payload)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(code)
        w.Write(response)
}