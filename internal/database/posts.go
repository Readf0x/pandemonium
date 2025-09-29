package database

import (
	"fmt"
	"slander/internal/systems"
	"strings"
	"time"
)

func PostExists(id systems.PostID) (exists bool, err error) {
	row := DB.QueryRow("select exists(select 1 from posts where id = ?)", id)
	err = row.Scan(&exists)
	return
}

func GetPost(id systems.PostID) (systems.Post, error) {
	row := DB.QueryRow("select post_type from posts where id = ?", id)
	var post_type string
	if err := row.Scan(&post_type); err != nil {
		return systems.Post{}, err
	}
	var (
		parent        int32
		owner         string
		original_body string
		body          string
		shares        int32
		timedate      string
	)
	if systems.PostTypeFromString(post_type) != systems.Original {
		row = DB.QueryRow("select parent, owner, original_body, body, shares, time from posts where id = ?", id)
		if err := row.Scan(&parent, &owner, &original_body, &body, &shares, &timedate); err != nil {
			return systems.Post{}, err
		}
	} else {
		row = DB.QueryRow("select owner, original_body, body, shares, time from posts where id = ?", id)
		if err := row.Scan(&owner, &original_body, &body, &shares, &timedate); err != nil {
			return systems.Post{}, err
		}
	}
	t, err := time.Parse(time.RFC3339, timedate)
	if err != nil {
		return systems.Post{}, err
	}
	id2 := systems.PostID(id)
	likes, err := GetLikeCount(id2)
	if err != nil {
		return systems.Post{}, err
	}
	return systems.Post{
		ID:           id2,
		Parent:       systems.PostID(parent),
		Owner:        systems.UserID(owner),
		OriginalBody: original_body,
		Body:         body,
		Likes:        likes,
		Shares:       shares,
		Time:         t,
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
			"insert into posts (parent, owner, original_body, body, shares, time, post_type) values (?, ?, ?, ?, ?, ?, ?)",
			post.Parent,
			post.Owner,
			post.Body,
			post.Body,
			post.Shares,
			post.Time.Format(time.RFC3339),
			post_type.String(),
		)
	} else {
		_, err = DB.Exec(
			"insert into posts (owner, original_body, body, shares, time, post_type) values (?, ?, ?, ?, ?, ?)",
			post.Owner,
			post.Body,
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

func LikePost(user systems.UserID, post systems.PostID) (err error) {
	_, err = DB.Exec("insert into likes (user_id, post_id) values (?, ?)", user, post)
	return
}

func UnlikePost(user systems.UserID, post systems.PostID) (err error) {
	_, err = DB.Exec("delete from likes where user_id = ? and post_id = ?", user, post)
	return
}

func QueryLike(user systems.UserID, post systems.PostID) (liked bool, err error) {
	row := DB.QueryRow("select exists(select 1 from likes where user_id = ? and post_id = ?)", user, post)
	err = row.Scan(&liked)
	return
}

func DeletePost(id systems.PostID) (err error) {
	_, err = DB.Exec("delete from posts where id = ?", id)
	return
}

func IncrementShares(id systems.PostID) (err error) {
	_, err = DB.Exec("update posts set shares = shares + 1 where id = ?", id)
	return
}

func GetPage(size int, page int) (posts []systems.Post, err error) {
	rows, err := DB.Query("select id, owner, original_body, body, shares, time from posts where parent is null order by time desc limit ? offset ?", size, page * size)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id            systems.PostID
			owner         string
			original_body string
			body          string
			shares        int32
			timedate      string
		)
		if err = rows.Scan(&id, &owner, &original_body, &body, &shares, &timedate); err != nil {
			return
		}

		t, err := time.Parse(time.RFC3339, timedate)
		if err != nil {
			return posts, err
		}

		likes, err := GetLikeCount(id)
		if err != nil {
			return posts, err
		}

		posts = append(posts, systems.Post{
			ID:           systems.PostID(id),
			Owner:        systems.UserID(owner),
			OriginalBody: original_body,
			Body:         body,
			Likes:        likes,
			Shares:       shares,
			Time:         t,
		})
	}

	err = rows.Err()
	return
}

func EditPost(body string, id systems.PostID) (err error) {
	row := DB.QueryRow("select body from posts where id = ?", id)
	var current string
	err = row.Scan(&current)
	if err != nil {
		return
	}
	curw := strings.Split(current, " ")
	words := strings.Split(body, " ")
	changeCnt := 0
	if len(curw) != len(words) {
		err = fmt.Errorf("Incorrect word count!")
		return
	}
	for i, w := range words {
		if curw[i] != w {
			changeCnt++
		}
	}
	if changeCnt > 1 {
		err = fmt.Errorf("Too many changes!")
		return
	}
	_, err = DB.Exec("update posts set body = ? where id = ?", body, id)
	return
}
