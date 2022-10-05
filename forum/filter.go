package forum

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func allForumUnames() []string {
	var allUnames []string
	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var uname string
		rows.Scan(&uname)
		allUnames = append(allUnames, uname)
	}
	return allUnames
}

func AllForumUsers() []user {
	fmt.Println("ALLFORUMFIRSTLINE")
	var AllUsers []user
	rows, err := db.Query("SELECT * FROM users GROUP BY username;")
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
		os.Exit(0)
	}
	defer rows.Close()
	for rows.Next() {
		var usr user
		rows.Scan(&(usr.Username), &(usr.Image), &(usr.Email), &(usr.Password), &(usr.Access), &(usr.LoggedIn), &(usr.LikedPost), &(usr.DislikedPost), &(usr.LikedComments2), &(usr.DislikedComments2), &(usr.Notifyview), &(usr.Notifymsg), &(usr.LikedComments))
		AllUsers = append(AllUsers, usr)

	}
	return AllUsers
}

func NotifForumUsers() []user {
	fmt.Println("ALLFORUMFIRSTLINE")
	var AllUsers []user
	rows, err := db.Query("SELECT notifymsg FROM users GROUP BY username;")
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
		os.Exit(0)
	}
	defer rows.Close()
	for rows.Next() {
		var usr user
		rows.Scan(&(usr.Notifymsg))
		AllUsers = append(AllUsers, usr)

	}
	return AllUsers
}

func filCatDisplayPostsAndComments(filCat string) []post {
	var pos []post
	// fmt.Printf("filCat is %s\n", filCat)
	filCat = "%(" + filCat + ")%"
	rows, err := db.Query("SELECT * FROM posts WHERE category LIKE ?;", filCat)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var po post
		rows.Scan(&(po.PostID), &(po.Author), &(po.Image), &(po.Title), &(po.Content), &(po.Category), &(po.PostTime), &(po.Likes), &(po.Dislikes), &(po.IPs))
		po.PostTimeStr = po.PostTime.Format("Mon 02-01-2006 15:04:05")
		// fmt.Printf("Display Post: %d, by %s, title: %s, content: %s, in %s, at %v, with %d likes, and %d dislikes\n", po.PostID, po.Author, po.Title, po.Content, po.Category, po.PostTimeStr, po.Likes, po.Dislikes)
		po.Category = strings.Trim(po.Category, "(")
		po.Category = strings.Trim(po.Category, ")")
		po.Comments = displayComments(po.PostID)
		pos = append(pos, po)
	}
	return pos
}

func filAuthorDisplayPostsAndComments(filAuthor string) []post {
	var pos []post
	// fmt.Printf("filAuthor is %s\n", filAuthor)
	rows, err := db.Query("SELECT * FROM posts WHERE author LIKE ?;", filAuthor)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var po post
		rows.Scan(&(po.PostID), &(po.Author), &(po.Image), &(po.Title), &(po.Content), &(po.Category), &(po.PostTime), &(po.Likes), &(po.Dislikes), &(po.IPs))
		po.PostTimeStr = po.PostTime.Format("Mon 02-01-2006 15:04:05")
		// fmt.Printf("Display Post: %d, by %s, title: %s, content: %s, in %s, at %v, with %d likes, and %d dislikes\n", po.PostID, po.Author, po.Title, po.Content, po.Category, po.PostTimeStr, po.Likes, po.Dislikes)
		po.Comments = displayComments(po.PostID)
		po.Category = strings.Trim(po.Category, "(")
		po.Category = strings.Trim(po.Category, ")")
		pos = append(pos, po)
	}
	return pos
}

func filLikedDisplayPostsAndComments(curUser user) []post {
	var pos []post
	pos2 := displayPostsAndComments()
	AllLikes := CountLikesByUser(curUser, "l")
	for i := 0; i < len(pos2); i++ {
		for k := 0; k < len(AllLikes); k++ {
			if pos2[i].PostID == AllLikes[k] {
				pos = append(pos, pos2[i])
			}
		}
	}
	return pos
}

func FindPostByPostId(postId int) post {
	var pos post
	po := displayPostsAndComments()
	for i := 0; i < len(po); i++ {
		if po[i].PostID == postId {
			pos = po[i]
		}
	}
	return pos
}

func FindComByComId(comId int) comment {
	var com comment
	co := displayComs()
	for i := 0; i < len(co); i++ {
		if co[i].CommentID == comId {
			com = co[i]
		}
	}
	return com
}
