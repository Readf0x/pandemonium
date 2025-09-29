package database

import (
	"crypto/sha256"
	"slander/internal/systems"
)

func CalcPassword(username string, password string) []byte {
	h := sha256.Sum256([]byte(username+password))
	return []byte(h[:])
}

func CheckPassword(id systems.UserID, hash []byte) (b bool, err error) {
	row := DB.QueryRow("select exists(select 1 from passwords where user_id = ? and hash = ?)", id, hash)
	err = row.Scan(&b)
	return
}

func SavePassword(hash []byte) (err error) {
	_, err = DB.Exec("insert into passwords (hash) values (?)", hash)
	return
}

