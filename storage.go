package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

func FreshNewDb() *Database {
	empty := []byte("{}")
	err := os.WriteFile("data.json", empty, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	db := Database{}
	if err := db.loadDatabase(); err != nil {
		log.Fatalln(err)
	}

	return &db
}

func (d *Database) loadDatabase() error {
	f, err := os.Open("data.json")
	if err != nil {
		return errors.New("failed to open database file to read")
	}

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&d); err != nil {
		log.Fatal(err)
		return errors.New("failed to read database")
	}

	id := len(d.Chirps)
	if id != 0 {
		d.latestChirpId = d.Chirps[len(d.Chirps)-1].Id
	}
	id = len(d.Users)
	if id != 0 {
		d.latestUserId = d.Users[len(d.Users)-1].Id
	}

	return nil
}

func (d *Database) sync() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	data, err := json.Marshal(d)
	if err != nil {
		return errors.New("failed to marshell database")
	}

	if err := os.WriteFile("data.json", data, 0644); err != nil {
		return errors.New("failed to write to disc")
	}

	return nil
}

func (d *Database) storeChirp(c Chirp) (Chirp, error) {
	d.mu.Lock()
	d.mu.Unlock()

	if c.Id == 0 {
		d.latestChirpId++
		c.Id = d.latestChirpId
		d.Chirps = append(d.Chirps, c)
	} else {
		d.Chirps[c.Id-1] = c
	}

	return c, nil
}

func (d *Database) storeUser(u User) (User, error) {
	d.mu.Lock()
	d.mu.Unlock()

	if u.Id == 0 {
		d.latestUserId++
		u.Id = d.latestUserId
		d.Users = append(d.Users, u)
	} else {
		d.Users[u.Id-1] = u
	}

	return u, nil
}

type Database struct {
	Chirps        []Chirp `json:"chirps"`
	latestChirpId int
	Users         []User `json:"users"`
	latestUserId  int
	mu            sync.Mutex
}
