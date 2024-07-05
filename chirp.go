package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func getProfanityWords() []string {
	return []string{"kerfuffle", "sharbert", "fornax"}
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" || !strings.Contains(bearer, "Bearer ") {
		resp := errorResponse{"Authorization token is missing"}
		dat, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed marshalling json error response")
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(dat)
		return
	}

	bearer = strings.Split(bearer, " ")[1]
	token, err := jwt.ParseWithClaims(
		bearer,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		},
	)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		resp := errorResponse{"Failed to get ID from jwt"}
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

	userID, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	var user *User
	for _, userInDb := range cfg.database.Users {
		if userInDb.Id == userID {
			user = &userInDb
		}
	}

	if user == nil {
		w.WriteHeader(401)
		return
	}

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

	chirp := Chirp{Body: strings.Join(parts, " "), AuthorId: user.Id}
	chirp, err = cfg.database.storeChirp(chirp)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to store chirp in database")
		w.WriteHeader(500)
		return
	}

	dat, err := json.Marshal(chirp)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed marshalling json error response")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	queryAuthor := r.URL.Query().Get("author_id")
	var authorId *int
	if queryAuthor != "" {
		id, err := strconv.Atoi(queryAuthor)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		authorId = &id
	}

	chirpMap := []Chirp{}
	for _, chirp := range cfg.database.Chirps {
		if authorId == nil || *authorId == chirp.AuthorId {
			chirpMap = append(chirpMap, chirp)
		}
	}

	data, err := json.Marshal(&chirpMap)
	if err != nil {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" || !strings.Contains(bearer, "Bearer ") {
		resp := errorResponse{"Authorization token is missing"}
		dat, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed marshalling json error response")
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(dat)
		return
	}

	bearer = strings.Split(bearer, " ")[1]
	token, err := jwt.ParseWithClaims(
		bearer,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		},
	)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		resp := errorResponse{"Failed to get ID from jwt"}
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

	userID, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	var user *User
	for _, userInDb := range cfg.database.Users {
		if userInDb.Id == userID {
			user = &userInDb
		}
	}

	if user == nil {
		w.WriteHeader(401)
		return
	}

	lookingFor, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		resp := errorResponse{"ID param is not a valid number"}
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

	var chirp *Chirp

	for _, chirpInDb := range cfg.database.Chirps {
		if chirpInDb.Id == lookingFor {
			chirp = &chirpInDb
		}
	}

	if chirp == nil {
		w.WriteHeader(404)
		return

	}

	if chirp.AuthorId != user.Id {
		w.WriteHeader(403)
		return
	}

	cfg.database.deleteChirp(*chirp)
	w.WriteHeader(204)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	lookingFor, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		resp := errorResponse{"ID param is not a valid number"}
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

	var chirp *Chirp

	for _, chirpInDb := range cfg.database.Chirps {
		if chirpInDb.Id == lookingFor {
			chirp = &chirpInDb
		}
	}

	if chirp == nil {
		resp := errorResponse{"Chirp does not exist"}
		dat, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed marshalling json error response")
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(dat)
		return

	}

	data, err := json.Marshal(chirp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

type errorResponse struct {
	Error string `json:"error"`
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}
