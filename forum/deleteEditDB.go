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
