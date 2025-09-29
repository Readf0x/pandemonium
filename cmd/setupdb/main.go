package main

import (
	"fmt"
	"os"
	"slander/internal/database"
)

func main() {
	if _, err := os.Stat("app.db"); err == nil {
		fmt.Println("'app.db' already exists!")
		fmt.Print("Are you sure? y/N: ")
		var c string
		fmt.Scan(&c)
		if c[0] != 'y' && c[0] != 'Y' {
			os.Exit(0)
		} else {
			os.Remove("app.db")
			database.DB = database.OpenDB()
		}
	}
	r := database.CalcPassword("readf0x", "asdf")
	fmt.Println(r)
	x := database.CalcPassword("xkgjl0d", "asdf")
	fmt.Println(x)
	_, err := database.DB.Exec(`
		insert into users (id, display_name, picture) values ("readf0x", "Tony Key Fobs", "NULL");
		insert into tokens (user_id, token, expiry) values ("readf0x", "debug", "3000-09-29T01:24:58.000Z");
		insert into passwords (user_id, hash) values ("readf0x", ?);
		insert into users (id, display_name, picture) values ("xkgjl0d", "Meow Mix Motherfucker", "NULL");
		insert into tokens (user_id, token, expiry) values ("xkgjl0d", "debug", "3000-09-29T01:24:58.000Z");
		insert into passwords (user_id, hash) values ("xkgjl0d", ?);
	`, []byte(r[:]), []byte(x[:]))
	if err != nil {
		fmt.Println(err)
		os.Remove("./app.db")
		os.Exit(1)
	}
}

