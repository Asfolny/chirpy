package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	database       Database
	jwtSecret      string
	polkaKey       string
}

func main() {
	mux := http.NewServeMux()
	godotenv.Load()
	apiCfg := apiConfig{
		database:  *FreshNewDb(),
		jwtSecret: os.Getenv("JWT_SECRET"),
		polkaKey:  os.Getenv("POLKA_KEY"),
	}

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
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
