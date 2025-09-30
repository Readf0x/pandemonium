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
	"slander/internal/routes"
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
		Error_Handler(err, w, r)
		return
	}
}

func Error_Handler(err error, w http.ResponseWriter, r *http.Request) {
	// log.Printf("Problem: %v\n%s", err, debug.Stack())
	w.Header().Set("HX-Retarget", "#errorBanner")
	w.Header().Set("HX-Reswap", "outerHTML")
	components.ErrorBanner(err).Render(r.Context(), w)
}

func GenericRoute_Handler(page templ.Component) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		page.Render(r.Context(), w)
	}
}

func Route_Handler(title string, page string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		routes.Page(title, page).Render(r.Context(), w)
	}
}

func Page_Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	userID, err := ValidateUser(r)
	if err != nil {
		w.Header().Set("HX-Redirect", "/login")
		return
	}
	user, err := database.GetUser(userID)
	if err != nil {
		Error_Handler(err, w, r)
		return
	}
	switch path[1] {
	case "home":
		routes.Home(user).Render(r.Context(), w)
	}
}

func Root_Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		GenericRoute_Handler(routes.Root())(w, r)
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
		userID, postID, err := CreatePost(r)
		if err != nil {
			components.PostSendButton(err, true).Render(r.Context(), w)
			Error_Handler(err, w, r)
			return
		}
		user, err := database.GetUser(userID)
		if err != nil {
			components.PostSendButton(err, true).Render(r.Context(), w)
			Error_Handler(err, w, r)
			return
		}
		post, err := database.GetPost(postID)
		if err != nil {
			components.PostSendButton(err, true).Render(r.Context(), w)
			Error_Handler(err, w, r)
			return
		}
		w.Header().Set("HX-Retarget", "#newPost")
		w.Header().Set("HX-Reswap", "outerHTML")
		components.PostSendButton(systems.PostSent200, true).Render(r.Context(), w)
		components.ResetPostSendButton(systems.PostSent200, true).Render(r.Context(), w)
		components.NewPost(components.Post(user, post)).Render(r.Context(), w)
	case "validate":
		_, err := ValidateUser(r)
		if err != nil {
			if r.FormValue("redirect") != "false" {
				w.Header().Set("HX-Redirect", "/login")
			}
			return
		}
		w.Header().Set("HX-Redirect", "/home")
	case "signup":
		r.ParseForm()
		username := r.FormValue("username")
		database.CreateUser(systems.User{
			ID: systems.UserID(username),
			DisplayName: r.FormValue("display_name"),
			Picture: "NULL",
		})
		database.SavePassword(
			database.CalcPassword(username, r.FormValue("password")),
		)
	case "login":
		r.ParseForm()
		username := r.FormValue("username")
		id := systems.UserID(username)
		hash := database.CalcPassword(username, r.FormValue("password"))
		b, err := database.CheckPassword(id, hash)
		if !b || err != nil {
			components.LoginButton(err).Render(r.Context(), w)
		} else {
			token, err := database.GenerateToken(id)
			if err != nil {
				log.Println(err)
				return
			}
			w.Header().Set("HX-Retarget", "#returnScript")
			w.Header().Set("HX-Reswap", "outerHTML")
			components.LogIn(id, token).Render(r.Context(), w)
		}
	}
}

func ValidateUser(r *http.Request) (id systems.UserID, err error) {
	if r.Form == nil {
		err = r.ParseForm()
		if err != nil {
			return
		}
	}
	id = systems.UserID(r.FormValue("user"))
	token := systems.SessionToken(r.FormValue("token"))
	tokens, err := database.GetUserTokens(id)
	if !slices.Contains(tokens, token) {
		err = fmt.Errorf("Invalid token!")
	}
	return
}

func CreatePost(r *http.Request) (userID systems.UserID, postID systems.PostID, err error) {
	userID, err = ValidateUser(r)
	body := r.FormValue("body")
	if systems.SpeechFilter.IsProfane(body) {
		err = fmt.Errorf("This content may violate our terms of service.")
		return
	}
	if len(body) > 128 {
		err = fmt.Errorf("Too many characters")
		return
	}
	if strings.Trim(body, " \n") == "" {
		err = fmt.Errorf("Empty post body")
		return
	}
	post_type := systems.PostTypeFromString(r.FormValue("type"))
	if post_type == systems.Original {
		postID, err = database.CreatePost(systems.Post{
			Owner:    userID,
			Body:     body,
			Time:     time.Now(),
			PostType: systems.Original,
		})
	} else {
		p := r.FormValue("parent")
		parent, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return userID, postID, err
		}
		postID, err = database.CreatePost(systems.Post{
			Owner:    userID,
			Body:     body,
			Time:     time.Now(),
			Parent:   systems.PostID(parent),
			PostType: post_type,
		})
	}
	return
}

func PostAction(path []string, w http.ResponseWriter, r *http.Request) {
	raw, err := strconv.ParseInt(path[2], 10, 32)
	if err != nil {
		Error_Handler(err, w, r)
		return
	}
	postID := systems.PostID(raw)
	userID, err := ValidateUser(r)
	if err != nil {
		Error_Handler(err, w, r)
		return
	}
	switch path[3] {
	case "comment":
		user, err := database.GetUser(userID)
		if err != nil {
			Error_Handler(err, w, r)
			return
		}
		post, err := database.GetPost(postID)
		if err != nil {
			Error_Handler(err, w, r)
			return
		}
		components.FloatReply(post, user).Render(r.Context(), w)
	case "repost":
	case "like":
		err := database.LikePost(userID, postID)
		if err != nil {
			Error_Handler(err, w, r)
			return
		}
		count, err := database.GetLikeCount(postID)
		if err != nil {
			Error_Handler(err, w, r)
			return
		}
		w.Header().Set("HX-Reswap", "outerHTML")
		components.PostActionButton(fmt.Sprintf("/htmx/postaction/%d/%s", postID, "unlike"), "heart-fill", count, "text-red-400 dark:text-red-500").Render(r.Context(), w)
	case "unlike":
		err := database.UnlikePost(userID, postID)
		if err != nil {
			Error_Handler(err, w, r)
			return
		}
		count, err := database.GetLikeCount(postID)
		if err != nil {
			Error_Handler(err, w, r)
			return
		}
		w.Header().Set("HX-Reswap", "outerHTML")
		components.PostActionButton(fmt.Sprintf("/htmx/postaction/%d/%s", postID, "like"), "heart", count, "hover:text-red-400 dark:hover:text-red-500").Render(r.Context(), w)
	case "delete":
		row := database.DB.QueryRow("select owner from posts where id = ?", postID)
		var ownerID systems.UserID
		err = row.Scan(&ownerID)
		if err != nil {
			Error_Handler(err, w, r)
			return
		}
		if userID != ownerID {
			Error_Handler(fmt.Errorf("Invalid permissions!"), w, r)
			return
		}
		err = database.DeletePost(postID)
		if err != nil {
			Error_Handler(err, w, r)
			return
		}
		w.Header().Set("HX-Retarget", fmt.Sprintf("#post%d", postID))
		w.Header().Set("HX-Reswap", "outerHTML")
		fmt.Fprint(w, "")
	case "resetbody":
		w.Header().Set("HX-Retarget", fmt.Sprintf("#post%d .body", postID))
		w.Header().Set("HX-Reswap", "innerHTML")
		fmt.Fprint(w, r.FormValue("body"))
	case "editmode":
		w.Header().Set("HX-Retarget", fmt.Sprintf("#post%d .body", postID))
		w.Header().Set("HX-Reswap", "innerHTML")
		components.EditModeBody(r.FormValue("body"), postID).Render(r.Context(), w)
	case "edit":
		body := r.FormValue("body")
		if systems.SpeechFilter.IsProfane(body) {
			Error_Handler(fmt.Errorf("This content may violate our terms of service."), w, r)
			return
		}
		err := database.EditPost(body, postID)
		if err != nil {
			Error_Handler(err, w, r)
			return
		}
		w.Header().Set("HX-Retarget", "#editBody")
		w.Header().Set("HX-Reswap", "outerHTML")
		fmt.Fprint(w, body)
	}
}

func Reset_Handler(path []string, w http.ResponseWriter, r *http.Request) {
	switch path[2] {
	case "postSendButton":
		components.PostSendButton(nil, true).Render(r.Context(), w)
	case "errorBanner":
		components.ErrorBanner(nil).Render(r.Context(), w)
	case "floatUI":
		fmt.Fprint(w, `<div id="postInput" hx-swap-oob="true"></div>`)
	}
}
