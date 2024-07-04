package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	mux := http.NewServeMux()

	apiCfg := apiConfig{}

	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))),
	)
	mux.HandleFunc("GET /api/healthz", apiCfg.handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerGetMetrics)
	mux.HandleFunc("/api/reset", apiCfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handleChirpValidate)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
