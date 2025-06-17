package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
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