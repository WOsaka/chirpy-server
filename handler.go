package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	hits := cfg.fileserverHits.Load()
	htmlTempl := `
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`
	html := fmt.Sprintf(htmlTempl, hits)
	w.Write([]byte(html))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Hits counter reset"))
}

func validationHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	// type validationResp struct {
	// 	Error        string `json:"error,omitempty"`
	// 	Cleaned_body string `json:"cleaned_body,omitempty"`
	// }

	// var respBody validationResp
	chirp := params.Body
	if len(chirp) > 140 {
		// respBody.Error = "Chirp is to long"
		// w.WriteHeader(400)
		respondWithError(w, 400, "Chirp is to long")
	} else {
		chirp := replaceProfane(chirp)
		respBody := struct {
			Cleaned_body string `json:"cleaned_body"`
		}{
			Cleaned_body: chirp,
		}
		respondWithJSON(w, 200, respBody)
	}

	// dat, err := json.Marshal(respBody)
	// if err != nil {
	// 	log.Printf("Error marshalling JSON: %s", err)
	// 	w.WriteHeader(500)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	// w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type Error struct {
		Error string `json:"error"`
	}

	var respBody Error
	respBody.Error = msg
	w.WriteHeader(code)

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func replaceProfane(sentence string) string {
	wordsToCheck := []string{"Kerfuffle", "Sharbert", "Fornax"}

	for _, word := range wordsToCheck {
		lowerWord := strings.ToLower(word)
		if strings.Contains(sentence, word) || strings.Contains(sentence, lowerWord) {
			sentence = strings.ReplaceAll(sentence, word, "****")
			sentence = strings.ReplaceAll(sentence, lowerWord, "****")
		}
	}

	return sentence
}
