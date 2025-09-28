package handler

import (
	"net/http"
	"slander/internal/routes"
)

func Route(mux *http.ServeMux) {
	mux.HandleFunc("/", Root_Handler)
	mux.HandleFunc("/assets/", Asset_Handler)
	mux.HandleFunc("/page/", Page_Handler)
	mux.HandleFunc("/home", Route_Handler("Home", "home"))
	mux.HandleFunc("/login", GenericRoute_Handler(routes.Login()))
	mux.HandleFunc("/htmx/", HTMX_Handler)
}
