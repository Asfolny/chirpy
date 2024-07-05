package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	passHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to make bcrypt hash: %s\n", err)

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

	user := User{Email: params.Email, Password: string(passHash)}
	user, err = cfg.database.storeUser(user)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to store user in database")
		w.WriteHeader(500)
		return
	}

	dat, err := json.Marshal(user)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed marshalling json error response")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	userMap := make([]User, len(cfg.database.Users))
	for i, user := range cfg.database.Users {
		userMap[i] = user
	}

	data, err := json.Marshal(&userMap)
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

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	lookingFor, err := strconv.Atoi(r.PathValue("userID"))
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

	var user *User

	for _, userInDb := range cfg.database.Users {
		if userInDb.Id == lookingFor {
			user = &userInDb
		}
	}

	if user == nil {
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

	data, err := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	var user *User

	for _, userInDb := range cfg.database.Users {
		if userInDb.Email == params.Email {
			user = &userInDb
		}
	}

	if user == nil {
		resp := errorResponse{"User does not exist"}
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

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		w.WriteHeader(401)
		return
	}

	type userResponse struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	userResp := userResponse{user.Id, user.Email}
	data, err := json.Marshal(&userResp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
