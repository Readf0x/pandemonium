package handler

import (
	"net/http"
	"slander/internal/routes"
)

func Route(mux *http.ServeMux) {
	mux.HandleFunc("/", Root_Handler)
	mux.HandleFunc("/assets/", Asset_Handler)
	mux.HandleFunc("/home", Page_Handler(routes.Home()))
	mux.HandleFunc("/htmx/", HTMX_Handler)
}
