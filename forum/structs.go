package forum

import (
	"html/template"
	"time"
)

type comment struct {
	CommentID      int
	Author         string
	PostID         int
	Content        string
	CommentTime    time.Time
	CommentTimeStr string
	Likes          int
	Dislikes       int
	LikedByCur     bool
	DislikedByCur  bool
	URL            string
	Deleted bool
}

type post struct {
	PostID        int
	Author        string // author
	Image         string
	Title         string
	Content       string
	Category      string
	PostTime      time.Time
	PostTimeStr   string
	Likes         int
	Dislikes      int
	Comments      []comment
	IPs           string
	View          int
	LikedByCur    bool
	DislikedByCur bool
	URL           string
	Deleted bool
}

type user struct {
	Username          string
	Email             string
	Access            int // 0 means no access, not logged in
	LoggedIn          bool
	Image             string
	Posts             []post
	Comments          []comment
	LikedPost         string
	DislikedPost      string
	LikedComments2    string
	DislikedComments2 string
	LikedComments     []comment
	Password          string
	Notifyview        string
	Notifymsg         string
	NotifMessageShow  map[string]template.URL
}

type Activity struct {
	Username   string
	Likes      string
	Dislikes   string
	PostID     string
	CommentID  string
	LikesCom   string
	DlikesCom  string
	Notifymsg  string
	Post       []post
	Com        []comment
	LikedPost  []post
	DlikedPost []post
	LikedCom   []comment
	DlikedCom  []comment
}
