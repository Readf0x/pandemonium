package main

import (
	"log"
	"net/http"
	"slander/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	handler.Route(mux)
	log.Fatal(http.ListenAndServe(":3000", mux))
}
