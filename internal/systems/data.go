package systems

import (
	"fmt"
	"time"
)

var PostSent200 = fmt.Errorf("PostSent200")

type UserID string
type SessionToken string

type User struct {
	ID          UserID
	DisplayName string
	Picture     string
}

type PostID int64
type PostType int
const (
	Original PostType = iota
	Comment
	Repost
)
var (
	postTypeString = map[PostType]string{
		Original: "original",
		Comment:  "comment",
		Repost:   "repost",
	}
)
func (p PostType) String() string {
	return postTypeString[p]
}

type Post struct {
	ID       PostID
	Parent   PostID
	Owner    UserID
	Body     string
	Likes    int32
	Shares   int32
	Time     time.Time
}
