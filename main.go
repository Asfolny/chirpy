package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
