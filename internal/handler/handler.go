package handler

import (
	"log"
	"net/http"
	"strings"
	"os"
	"io"
	"mime"
	"path/filepath"

	"github.com/a-h/templ"
)

func Asset_Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	asset, err := os.Open(path)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	defer asset.Close()
	ext := strings.ToLower(filepath.Ext(path))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		w.Header().Set("Content-Type", "application/octect-stream")
	} else {
		w.Header().Set("Content-Type", mimeType)
	}
	log.Printf("Serving asset '%s' with Content-Type: %s", path, w.Header().Get("Content-Type"))
	_, err = io.Copy(w, asset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func Page_Handler(page templ.Component) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		page.Render(r.Context(), w)
	}
}

func Root_Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.NotFound(w, r)
	}
}

func HTMX_Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Write([]byte("lorem ipsum dolor sit amet"))
}
