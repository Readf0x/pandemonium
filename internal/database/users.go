package database

import "slander/internal/systems"

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
		if err = rows.Scan(&token); err != nil {
			return
		}
		tokens = append(tokens, systems.SessionToken(token))
	}

	err = rows.Err()
	return
}

func CreateUser(user systems.User) (err error) {
	_, err = DB.Exec("insert into users (id, display_name, picture) values (?, ?, ?)", user.ID, user.DisplayName, user.Picture)
	return
}
