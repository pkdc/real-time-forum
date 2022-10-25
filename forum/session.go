package forum

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func genCookie(conn *websocket.Conn, uid int) SessionCookie {
	// assign a cookie
	sid := uuid.NewV4()
	fmt.Printf("login sid: %s for uid %d\n", sid, uid)

	var cookieResp SessionCookie

	cookieResp.Uid = uid
	cookieResp.Sid = sid.String()
	cookieResp.MaxAge = 1800

	// http.SetCookie(w, &http.Cookie{
	// 	Name:   "session",
	// 	Value:  sid.String(),
	// 	MaxAge: 900, // 15mins
	// })

	// allow each user to have only one opened session
	var loggedInUid int
	rows, err := db.Query("SELECT userID FROM sessions WHERE userID = ?;", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&loggedInUid)
	}
	// if the uid can be found in session table, remove that row (should only have 1 row)
	if loggedInUid != 0 {
		stmt, err := db.Prepare("DELETE FROM sessions WHERE userID = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(loggedInUid)
	}

	// // forumUser.Username = nicknameEmailDB
	// // forumUser.LoggedIn = true
	// // fmt.Printf("%s forum User Login\n", forumUser.Username)

	// update the user's login status
	stmt, err := db.Prepare("UPDATE users SET loggedIn = ? WHERE userID = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(true, uid)

	// // insert a record into session table
	stmt, err = db.Prepare("INSERT INTO sessions (sessionID, userID) VALUES (?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(sid.String(), uid)

	//test
	// var whichUser int
	// var logInOrNot bool
	// rows, err = db.Query("SELECT userID, loggedIn FROM users WHERE userID = ?;", uid)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	rows.Scan(&whichUser, &logInOrNot)
	// }
	// fmt.Printf("login user: %d, login status: %v\n", whichUser, logInOrNot)

	return cookieResp
}
