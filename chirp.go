package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
)

func getProfanityWords() []string {
	return []string{"kerfuffle", "sharbert", "fornax"}
}

func handleChirpValidate(w http.ResponseWriter, r *http.Request) {
	maxLen := 140

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	if err := decoder.Decode(&params); err != nil {
		fmt.Fprintf(os.Stderr, "error decoding parameters: %s\n", err)

		resp := errorResponse{"Something went wrong"}
		dat, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed marshalling json error response")
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	if len(params.Body) > maxLen {
		resp := errorResponse{"Chirp is too long"}
		dat, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed marshalling json error response")
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	parts := strings.Split(params.Body, " ")
	profanity := getProfanityWords()

	for i, word := range parts {
		if slices.Contains(profanity, strings.ToLower(word)) {
			parts[i] = "****"
		}
	}

	resp := validResponse{strings.Join(parts, " ")}
	dat, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed marshalling json error response")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

type validResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorResponse struct {
	Error string `json:"error"`
}