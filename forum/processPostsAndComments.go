package forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func processPost(r *http.Request, curUser user) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	idNumOfLikesStr := r.PostForm.Get("po-like")
	idNumOfDislikesStr := r.PostForm.Get("po-dislike")
	postTitle := r.PostForm.Get("postTitle")

	if idNumOfLikesStr != "" {
		fmt.Printf("current User username when liking post: %s\n", curUser.Username)
		idNumOfLikesStrSlice := strings.Split(idNumOfLikesStr, "-")
		updatePostID := idNumOfLikesStrSlice[0]
		posID, err := strconv.Atoi(idNumOfLikesStrSlice[0])
		if err != nil {
			log.Fatal(err)
		}
		postIDwithSep := curUser.LikedPost + "-" + updatePostID

		stmt2, err := db.Prepare("UPDATE users SET likedPosts = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt2.Close()
		stmt2.Exec(postIDwithSep, curUser.Username)
		curUser.DislikedPost = CheckLikesAndDislikes(curUser, posID, "d")
		stmt3, err := db.Prepare("UPDATE users SET dislikedPosts = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt3.Close()
		stmt3.Exec(curUser.DislikedPost, curUser.Username)
		url, author := findAuthor(posID)
		randomID := RandStringRunes(10)
		author.Notifymsg = author.Notifymsg + curUser.Username + " liked your post" + url + randomID + "#"
		// author.NotifView += randomID+ "#"
		stmt, err := db.Prepare("UPDATE users SET Notifymsg = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(author.Notifymsg, author.Username)
		// stmt4, err := db.Prepare("UPDATE users SET Notifyview = ?	WHERE username = ?;")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer stmt4.Close()
		// stmt4.Exec(author.NotifView, author.Username)

	} else if idNumOfDislikesStr != "" {
		fmt.Printf("current User username when disliking post: %s\n", curUser.Username)
		idNumOfDislikesStrSlice := strings.Split(idNumOfDislikesStr, "-")
		updatePostID := idNumOfDislikesStrSlice[0]
		posID, err := strconv.Atoi(idNumOfDislikesStrSlice[0])
		if err != nil {
			log.Fatal(err)
		}
		postIDwithSep := curUser.DislikedPost + "-" + updatePostID
		stmt2, err := db.Prepare("UPDATE users SET dislikedPosts = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt2.Close()
		stmt2.Exec(postIDwithSep, curUser.Username)
		curUser.LikedPost = CheckLikesAndDislikes(curUser, posID, "l")
		stmt3, err := db.Prepare("UPDATE users SET likedPosts = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt3.Close()
		stmt3.Exec(curUser.LikedPost, curUser.Username)
		url, author := findAuthor(posID)
		randomID := RandStringRunes(10)
		author.Notifymsg = author.Notifymsg + curUser.Username + " disliked your post" + url + randomID + "#"
		// author.NotifView += randomID+ "#"
		stmt, err := db.Prepare("UPDATE users SET Notifymsg = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(author.Notifymsg, author.Username)
	} else if postTitle != "" {
		fmt.Printf("curUser username when inserting new post: %s\n", curUser.Username)
		postCon := r.PostForm.Get("postContent")
		postCat := r.PostForm["postCat"]
		var ip string

		// // Insert the first cat
		// stmt, err := db.Prepare("INSERT INTO posts (author, image, title, content, category, postTime, likes, dislikes, ips) VALUES (?,?,?,?,?,?,?,?,?);")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer stmt.Close()
		// stmt.Exec(curUser.Username, curUser.Image, postTitle, postCon, postCat[0], time.Now(), 0, 0, ip)

		// // Insert other cats if any, with the prev postID
		// if len(postCat) > 1 {
		// 	// fmt.Println("More than 1 post cat")
		// 	var curPostId int
		// 	rows, err := db.Query("SELECT MAX(postID) from posts;")
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	defer rows.Close()

		// 	for rows.Next() {
		// 		rows.Scan(&curPostId)
		// 	}

		// 	for cat := 1; cat < len(postCat); cat++ {
		// 		fmt.Printf("current post id is %d, cureently inserting the %dth category\n", curPostId, cat)
		// 		stmt, err = db.Prepare("INSERT INTO posts (postID, author, image, title, content, category, postTime, likes, dislikes, ips) VALUES(?,?,?,?,?,?,?,?,?,?);")
		// 		if err != nil {
		// 			log.Fatal(err)
		// 		}
		// 		defer stmt.Close()
		// 		stmt.Exec(curPostId, curUser.Username, curUser.Image, postTitle, postCon, postCat[cat], time.Now(), 0, 0, ip)
		// 	}
		// }

		postCatStr := ""
		for i := 0; i < len(postCat); i++ {
			postCatStr += "(" + postCat[i] + ")"
		}
		postID := strconv.Itoa(len(displayPostsAndComments()) + 1)
		url := "postpage?postdetails=" + postID + "&postdetails=" + postTitle
		stmt, err := db.Prepare("INSERT INTO posts (author, image, title, content, category, postTime, likes, dislikes, ips,URL,deleted) VALUES(?,?,?,?,?,?,?,?,?,?,?);")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(curUser.Username, curUser.Image, postTitle, postCon, postCatStr, time.Now(), 0, 0, ip, url,false)

		// test
		// var pid int
		// var un string
		// var t string
		// var con string
		// var cat string
		// var pT time.Time
		// var l int
		// var d int

		// rows, err := db.Query("SELECT * FROM posts")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer rows.Close()

		// for rows.Next() {
		// 	rows.Scan(&pid, &un, &t, &con, &cat, &pT, &l, &d)
		// 	fmt.Printf("Post: %d, by %s, Title: %s, content: %s, in %s, at %v, with %d likes, and %d dislikes\n", pid, un, t, con, cat, pT, l, d)
		// }
	}
	return
}

func processComment(r *http.Request, curUser user) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	idNumOfLikesStr := r.PostForm.Get("com-like")
	idNumOfDislikesStr := r.PostForm.Get("com-dislike")
	fmt.Println("LIKES", idNumOfLikesStr)
	fmt.Println("DISLIKES", idNumOfDislikesStr)
	comCon := r.PostForm.Get("comment")
	if idNumOfLikesStr != "" {
		fmt.Printf("curUser username when liking comment: %s\n", curUser.Username)
		idNumOfLikesStrSlice := strings.Split(idNumOfLikesStr, "-")
		comID := idNumOfLikesStrSlice[1]
		fmt.Println("LIKES", idNumOfLikesStr, comID)
		comIDint, err := strconv.Atoi(comID)
		if err != nil {
			log.Fatal(err)
		}
		comIDwithSep := curUser.LikedComments2 + "-" + comID
		stmt2, err := db.Prepare("UPDATE users SET likedComments2 = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt2.Close()
		stmt2.Exec(comIDwithSep, curUser.Username)
		curUser.DislikedComments2 = CheckLikesAndDislikes(curUser, comIDint, "ComD")
		stmt3, err := db.Prepare("UPDATE users SET dislikedComments2 = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt3.Close()
		stmt3.Exec(curUser.DislikedComments2, curUser.Username)
		url, author := findCommentAuthor(comIDint)
		randomID := RandStringRunes(10)
		author.Notifymsg = author.Notifymsg + curUser.Username + " liked your comment" + url + randomID + "#"
		// author.NotifView += randomID+ "#"
		stmt, err := db.Prepare("UPDATE users SET Notifymsg = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(author.Notifymsg, author.Username)

	} else if idNumOfDislikesStr != "" {
		fmt.Printf("curUser username when disliking comment: %s\n", curUser.Username)
		idNumOfDislikesStrSlice := strings.Split(idNumOfDislikesStr, "-")
		comID := idNumOfDislikesStrSlice[1]
		comIDint, err := strconv.Atoi(comID)
		if err != nil {
			log.Fatal(err)
		}
		comIDwithSep := curUser.DislikedComments2 + "-" + comID
		stmt2, err := db.Prepare("UPDATE users SET dislikedComments2 = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt2.Close()
		stmt2.Exec(comIDwithSep, curUser.Username)
		curUser.LikedComments2 = CheckLikesAndDislikes(curUser, comIDint, "ComL")
		stmt3, err := db.Prepare("UPDATE users SET likedComments2 = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt3.Close()
		stmt3.Exec(curUser.LikedComments2, curUser.Username)
		url, author := findCommentAuthor(comIDint)
		randomID := RandStringRunes(10)
		author.Notifymsg = author.Notifymsg + curUser.Username + " disliked your comment" + url + randomID + "#"
		// author.NotifView += randomID+ "#"
		stmt, err := db.Prepare("UPDATE users SET Notifymsg = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(author.Notifymsg, author.Username)
	} else if comCon != "" {
		fmt.Printf("curUser username when inserting new comment: %s\n", curUser.Username)
		poId := r.PostForm.Get("post-id")
		posID, err := strconv.Atoi(poId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("comment: %s under %s\n", comCon, poId)
		po := displayPostsAndComments()
		var title string
		for i := 0; i < len(po); i++ {
			if po[i].PostID == posID {
				title = po[i].Title
			}
		}
		link := "postpage?postdetails=" + poId + "&postdetails=" + title
		stmt, err := db.Prepare("INSERT INTO comments (author, postID, content, commentTime, likes, dislikes, URL,deleted) VALUES (?,?,?,?,?,?,?,?);")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(curUser.Username, poId, comCon, time.Now(), 0, 0, link,false)
		url, author := findAuthor(posID)
		randomID := RandStringRunes(10)
		author.Notifymsg = author.Notifymsg + curUser.Username + " commented your post" + url + randomID + "#"
		// author.NotifView += randomID+ "#"
		stmt2, err := db.Prepare("UPDATE users SET Notifymsg = ?	WHERE username = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt2.Close()
		stmt2.Exec(author.Notifymsg, author.Username)
	}
	return
}

// func LikesAndDislikes(post,curUserPost []post )[]post {
// 	for i:=0 ; i < len(post); i++ {

// 	}
// }

func dup_count(list []string) map[string]int {
	duplicate_frequency := make(map[string]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := duplicate_frequency[item]

		if exist {
			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item] = 1 // else start counting from 1
		}
	}
	return duplicate_frequency
}
