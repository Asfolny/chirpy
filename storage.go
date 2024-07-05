package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

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

	d.latestChirpId++
	c.Id = d.latestChirpId
	d.Chirps = append(d.Chirps, c)
	return c, nil
}

type Database struct {
	Chirps        []Chirp `json:"chirps"`
	latestChirpId int
	mu            sync.Mutex
}
