package database

import (
	"database/sql"
	"log"

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
			id integer primary key autoincrement,
			post_type text not null,
			parent integer,
			owner text not null,
			original_body text not null,
			body text not null,
			shares integer not null default 0,
			time string not null,
			foreign key(parent) references posts(id),
			foreign key(owner) references users(id) on delete cascade
		);
		create table if not exists likes (
			user_id text,
			post_id integer,
			primary key (user_id, post_id),
			foreign key(user_id) references users(id) on delete cascade,
			foreign key(post_id) references posts(id) on delete cascade
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
