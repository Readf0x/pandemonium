package database

import (
	"database/sql"
	"log"
	"slander/internal/systems"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB = OpenDB()

func OpenDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatalf("Failed to open database!\n%s", err)
	}
	Setup(db)
	return db
}

func Setup(db *sql.DB) {
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;
		create table if not exists users (
			id text primary key,
			display_name text not null,
			picture text not null
		);
		create table if not exists posts (
			id integer primary key,
			post_type text not null,
			parent integer,
			owner text not null,
			body text not null,
			shares integer not null default 0,
			time string not null,
			foreign key(parent) references posts(id),
			foreign key(owner) references users(id)
		);
		create table if not exists likes (
			user_id text,
			post_id integer,
			primary key (user_id, post_id),
			foreign key(user_id) references users(id),
			foreign key(post_id) references posts(id)
		);
		create table if not exists tokens (
			user_id text,
			token text,
			primary key (user_id, token),
			foreign key(user_id) references users(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

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
	rows, err := DB.Query("select token from tokens where user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, systems.SessionToken(token))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func CreateUser(user systems.User) (err error) {
	_, err = DB.Exec("insert into users (id, display_name, picture) values (?, ?, ?)", user.ID, user.DisplayName, user.Picture)
	return
}

func PostExists(id systems.PostID) (exists bool, err error) {
	row := DB.QueryRow("select exists(select 1 from posts where id = ?)", id)
	err = row.Scan(&exists)
	return
}

func GetPost(id systems.PostID) (systems.Post, error) {
	row := DB.QueryRow("select parent, owner, body, shares, time from posts where id = ?", id)
	var (
		parent int32
		owner string
		body string
		shares int32
		timedate string
	)
	if err := row.Scan(&parent, &owner, &body, &shares, &timedate); err != nil {
		return systems.Post{}, err
	}
	t, err := time.Parse(time.RFC3339, timedate)
	if err != nil {
		return systems.Post{}, err
	}
	return systems.Post{
		ID: systems.PostID(id),
		Parent: systems.PostID(parent),
		Owner: systems.UserID(owner),
		Likes: 0,
		Shares: 0,
		Time: t,
	}, nil
}

func PostCount() (count int32, err error) {
	row := DB.QueryRow("select count(*) from posts")
	err = row.Scan(&count)
	return
}

func CreatePost(post systems.Post, post_type systems.PostType) (err error) {
	if post_type != systems.Original {
		_, err = DB.Exec(
			"insert into posts (parent, owner, body, shares, time, post_type) values (?, ?, ?, ?, ?, ?)",
			post.Parent,
			post.Owner,
			post.Body,
			post.Shares,
			post.Time.Format(time.RFC3339),
			post_type.String(),
		)
	} else {
		_, err = DB.Exec(
			"insert into posts (owner, body, shares, time, post_type) values (?, ?, ?, ?, ?)",
			post.Owner,
			post.Body,
			post.Shares,
			post.Time.Format(time.RFC3339),
			post_type.String(),
		)
	}
	return
}

func GetLikeCount(id systems.PostID) (count int32, err error) {
	row := DB.QueryRow("select count(*) from likes where post_id = ?", id)
	err = row.Scan(&count)
	return
}
