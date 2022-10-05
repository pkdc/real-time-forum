package forum

import (
	"log"
	"net/http"
)

func deleteUser(r *http.Request) {
	r.ParseForm()
	dUser := r.PostForm.Get("delete")
	stmt, err := db.Prepare("DELETE FROM users WHERE username = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(dUser)
}

func DeleteAllComments() {
	stmt, err := db.Prepare("DELETE FROM comments;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func DeleteAllPosts() {
	stmt, err := db.Prepare("DELETE FROM posts;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func DeleteAllUsers() {
	stmt, err := db.Prepare("DELETE FROM users;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func ClearComments() {
	stmt, err := db.Prepare("DROP TABLE IF EXISTS comments;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func ClearPosts() {
	stmt, err := db.Prepare("DROP TABLE IF EXISTS posts;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func ClearUsers() {
	stmt, err := db.Prepare("DROP TABLE IF EXISTS users;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func DeleteOnePost(postID int) {
	stmt, err := db.Prepare("DELETE FROM posts WHERE postID=?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(postID)
	co := displayComments(postID)
	for i := 0; i < len(co); i++ {
		DeleteOneCom(co[i].CommentID)
	}
}

func DeleteOneCom(comID int) {
	stmt, err := db.Prepare("DELETE FROM comments WHERE commentID=?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(comID)
}

func EditPost(newpos post) {
	stmt, err := db.Prepare("UPDATE posts SET title = ?, content= ?,postTime =?, URL=? WHERE postID = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(newpos.Title, newpos.Content, newpos.PostTime, newpos.URL, newpos.PostID)
}

func EditCom(newCom comment) {
	stmt, err := db.Prepare("UPDATE comments SET content= ?,commentTime =? WHERE commentID = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(newCom.Content, newCom.CommentTime, newCom.CommentID)
}
