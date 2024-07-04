package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits++
		writer.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(writer, request)
	})
}

func (cfg *apiConfig) handlerGetMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(
		writer,
		`
<html>
    <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
    </body>
</html>
`,
		cfg.fileserverHits,
	)
}

func (cfg *apiConfig) handlerResetMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
}
