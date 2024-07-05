package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type apiConfig struct {
	fileserverHits int
	database       Database
}

func main() {
	mux := http.NewServeMux()

	apiCfg := apiConfig{database: *FreshNewDb()}

	// Keep database up to date
	go func() {
		for {
			time.Sleep(10 * time.Second)
			fmt.Println("Syncing database...")
			err := apiCfg.database.sync()
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Sync success")
		}
	}()

	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))),
	)
	mux.HandleFunc("GET /api/healthz", apiCfg.handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerGetMetrics)
	mux.HandleFunc("/api/reset", apiCfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
