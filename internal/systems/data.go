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
	stringPostType = map[string]PostType{
		"original": Original,
		"comment":  Comment,
		"repost":   Repost,
	}
)
func (p PostType) String() string {
	return postTypeString[p]
}
func PostTypeFromString(str string) PostType {
	return stringPostType[str]
}

type Post struct {
	ID           PostID
	Parent       PostID
	Owner        UserID
	OriginalBody string
	Body         string
	Likes        int32
	Shares       int32
	Time         time.Time
}
