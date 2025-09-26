package main

import "slander/internal/database"

func main() {
	database.DB.Exec(`
		insert into users (id, display_name, picture) values ("readf0x", "Tony Key Fobs", "NULL");
		insert into tokens (user_id, token) values ("readf0x", "debug");
		insert into users (id, display_name, picture) values ("xkgjl0d", "Meow Mix Motherfucker", "NULL");
		insert into tokens (user_id, token) values ("xkgjl0d", "debug");
	`)
}

