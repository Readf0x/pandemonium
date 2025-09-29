package database

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"slander/internal/systems"
	"time"
)

func UserExists(id systems.UserID) (exists bool, err error) {
	row := DB.QueryRow("select exists(select 1 from users where id = ?)", id)
	err = row.Scan(&exists)
	return
}

func GetUser(id systems.UserID) (systems.User, error) {
	row := DB.QueryRow("select display_name, picture from users where id = ?", id)
	var display_name, picture string
	if err := row.Scan(&display_name, &picture); err != nil {
		return systems.User{}, err
	}
	if picture == "NULL" {
		picture = "/assets/user.png"
	}
	return systems.User{
		ID:          systems.UserID(id),
		DisplayName: display_name,
		Picture:     picture,
	}, nil
}

func GetUserTokens(id systems.UserID) (tokens []systems.SessionToken, err error) {
	rows, err := DB.Query("select token, expiry from tokens where user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		var expiry string
		if err = rows.Scan(&token, &expiry); err != nil {
			return
		}
		t, err := time.Parse(time.RFC3339, expiry)
		if err != nil {
			return tokens, err
		}
		if t.Before(time.Now()) {
			_, err = DB.Exec("delete from tokens where user_id = ? and token = ?", id, token)
		}
		tokens = append(tokens, systems.SessionToken(token))
	}

	err = rows.Err()
	return
}

func GenerateToken(id systems.UserID) (token systems.SessionToken, err error) {
	now := time.Now()
	expiry := now.Add(7*24*time.Hour)
	t := sha256.Sum224(fmt.Append([]byte{}, expiry.Unix() * now.Unix(), id))
	b := base64.StdEncoding.EncodeToString(t[:])
	token = systems.SessionToken(b)
	_, err = DB.Exec("insert into tokens (user_id, token, expiry) values (?, ?, ?)", id, token, expiry.Format(time.RFC3339))
	return
}

func CreateUser(user systems.User) (err error) {
	_, err = DB.Exec("insert into users (id, display_name, picture) values (?, ?, ?)", user.ID, user.DisplayName, user.Picture)
	return
}
