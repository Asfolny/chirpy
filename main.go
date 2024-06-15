package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits++
		writer.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(writer, request)
	})
}

func main() {
	mux := http.NewServeMux()

	apiCfg := apiConfig{}

	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))),
	)
	mux.HandleFunc("/healthz", health)
	mux.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(writer, "Hits: %v", apiCfg.fileserverHits)
	})
	mux.HandleFunc("/reset", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(http.StatusOK)
		apiCfg.fileserverHits = 0
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}

func health(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(writer, "OK")
}
