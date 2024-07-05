package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" || !strings.Contains(bearer, "ApiKey ") {
		resp := errorResponse{"Authorization is missing"}
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

	apiKey := strings.Split(bearer, " ")[1]
	if cfg.polkaKey != apiKey {
		w.WriteHeader(401)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
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

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	var user *User

	for _, userInDb := range cfg.database.Users {
		if userInDb.Id == params.Data.UserID {
			user = &userInDb
		}
	}

	if user == nil {
		w.WriteHeader(404)
		return
	}

	user.Red = true
	cfg.database.storeUser(*user)

	w.WriteHeader(204)
}
