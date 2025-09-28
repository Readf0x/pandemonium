package main

import (
	"fmt"
	"os"
	"slander/internal/database"
)

func main() {
	_, err := database.DB.Exec(`
		insert into users (id, display_name, picture) values ("readf0x", "Tony Key Fobs", "NULL");
		insert into tokens (user_id, token, expiry) values ("readf0x", "debug", "never");
		insert into users (id, display_name, picture) values ("xkgjl0d", "Meow Mix Motherfucker", "NULL");
		insert into tokens (user_id, token, expiry) values ("xkgjl0d", "debug", "never");
	`)
	if err != nil {
		fmt.Println(err)
		os.Remove("./app.db")
		os.Exit(1)
	}
}

