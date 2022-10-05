package forum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// var forumUser user

func regNewUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	uname := r.PostForm.Get("username")
	email := r.PostForm.Get("email")
	password := []byte(r.PostForm.Get("password"))
	image := r.PostForm.Get("ProfilePic")
	if image == "" {
		tpl, err := template.ParseFiles("./templates/regUnfilledFields.gohtml")
		if err != nil {
			http.Error(w, "Parsing Error", http.StatusInternalServerError)
		}
		tpl.Execute(w, nil)
		return
	}
	if strings.Trim(uname, " ") == "" {
		tpl, err := template.ParseFiles("./templates/regUnfilledFields.gohtml")
		if err != nil {
			http.Error(w, "Parsing Error", http.StatusInternalServerError)
		}
		tpl.Execute(w, nil)
		return
	}

	emailSlice := strings.Split(email, "@")
	if len(emailSlice) != 2 {
		tpl, err := template.ParseFiles("./templates/regInvalidEmail.gohtml")
		if err != nil {
			http.Error(w, "Parsing Error", http.StatusInternalServerError)
		}
		tpl.Execute(w, nil)
		return
	}

	if !strings.ContainsRune(emailSlice[1], 46) {
		tpl, err := template.ParseFiles("./templates/regInvalidEmail.gohtml")
		if err != nil {
			http.Error(w, "Parsing Error", http.StatusInternalServerError)
		}
		tpl.Execute(w, nil)
		return
	}

	// check if already exists
	rows, err := db.Query("SELECT username, email FROM users WHERE username = ? OR email = ?;", uname, email)
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		// http.Error(w, "username or email is already taken", http.StatusConflict)
		// http.Redirect(w, r, "/register", http.StatusSeeOther)

		// render faildeReg template
		tpl, err := template.ParseFiles("./templates/failedReg.gohtml")
		if err != nil {
			log.Fatal(err)
		}
		tpl.Execute(w, nil)
		return
	}

	hash, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO users (username, image, email, password, access, loggedIn, likedPosts, dislikedPosts, likedComments2, dislikedComments2, notifyView, notifyMsg, likedComments) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(uname, image, email, hash, 1, true, "", "", "", "", "", "", "")

	// test
	var u string
	var i string
	var e string
	var p []byte
	var a int
	var l bool

	rows, err = db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&u, &i, &e, &p, &a, &l)
	}
	fmt.Printf("uname: %s i: %s e: %s pw: %s, ac: %d, log: %t\n", u, i, e, p, a, l)

	// forumUser.Username = uname
	// forumUser.LoggedIn = true
	// forumUser.Access = 1
	// forumUser.Image = image

	sid := uuid.NewV4()
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  sid.String(),
		MaxAge: 1800,
	})
	fmt.Printf("reg sid: %s\n", sid)
	fmt.Printf("Reg and login as %s\n", uname)

	stmt, err = db.Prepare("INSERT INTO sessions (sessionID, username) VALUES (?,?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(sid.String(), uname)
}
