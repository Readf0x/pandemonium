package handler

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slander/internal/components"
	"slander/internal/database"
	"slander/internal/systems"
	"slices"
	"strconv"
	"strings"
	"time"

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
	// log.Printf("Serving asset '%s' with Content-Type: %s", path, w.Header().Get("Content-Type"))
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
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "No Access", http.StatusUnauthorized)
		return
	}
	path := strings.Split(r.URL.Path[1:], "/")
	switch path[1] {
	case "reset":
		Reset_Handler(path, w, r)
	case "postaction":
		PostAction(path, w, r)
	case "createpost":
		err := CreatePost(r)
		if err != nil {
			components.PostSendButton(err).Render(r.Context(), w)
			components.ErrorBanner(err).Render(r.Context(), w)
			return
		}
		components.PostSendButton(systems.PostSent200).Render(r.Context(), w)
		components.ResetPostSendButton(systems.PostSent200).Render(r.Context(), w)
	}
}

func CreatePost(r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	id := systems.UserID(r.FormValue("user"))
	token := systems.SessionToken(r.FormValue("token"))
	tokens, err := database.GetUserTokens(id)
	if err != nil {
		return err
	}
	if !slices.Contains(tokens, token) {
		return fmt.Errorf("Invalid token!")
	}
	err = database.CreatePost(systems.Post{
		Owner: id,
		Body:  r.FormValue("body"),
		Time:  time.Now(),
	}, systems.Original)
	return err
}

func PostAction(path []string, w http.ResponseWriter, r *http.Request) {
	raw, err := strconv.ParseInt(path[2], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	postID := systems.PostID(raw)
	switch path[3] {
	case "comment":
	case "like":
		count, err := database.GetLikeCount(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		components.PostActionButton(postID, "unlike", "heart-fill", count, "text-red-400 dark:text-red-500").Render(r.Context(), w)
	case "unlike":
		count, err := database.GetLikeCount(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		components.PostActionButton(postID, "like", "heart", count, "hover:text-red-400 dark:hover:text-red-500").Render(r.Context(), w)
	}
}

func Reset_Handler(path []string, w http.ResponseWriter, r *http.Request) {
	switch path[2] {
	case "postSendButton":
		components.PostSendButton(nil).Render(r.Context(), w)
	}
}
