package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerHealth(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(writer, "OK")
}
