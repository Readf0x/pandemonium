package main

import (
	"htmx-go-server/internal/handler"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	handler.Route(mux)
	log.Fatal(http.ListenAndServe(":3000", mux))
}

