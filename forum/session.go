package forum

import (
	"fmt"

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

	return cookieResp

	// http.SetCookie(w, &http.Cookie{
	// 	Name:   "session",
	// 	Value:  sid.String(),
	// 	MaxAge: 900, // 15mins
	// })

	// // allow each user to have only one opened session
	// var loggedInUname string
	// rows, err = db.Query("SELECT username FROM sessions WHERE username = ?;", nicknameEmailDB)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	rows.Scan(&loggedInUname)
	// }
	// // if the uname can be found in session table, remove that row (should only have 1 row)
	// if loggedInUname != "" {
	// 	stmt, err := db.Prepare("DELETE FROM sessions WHERE username = ?;")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer stmt.Close()
	// 	stmt.Exec(loggedInUname)
	// }

	// // forumUser.Username = nicknameEmailDB
	// // forumUser.LoggedIn = true
	// // fmt.Printf("%s forum User Login\n", forumUser.Username)

	// // update the user's login status
	// stmt, err := db.Prepare("UPDATE users SET loggedIn = ? WHERE userID = ?;")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// stmt.Exec(true, userIDDB)

	// // insert a record into session table
	// stmt, err = db.Prepare("INSERT INTO sessions (sessionID, userID) VALUES (?, ?);")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// stmt.Exec(sid.String(), nicknameEmailDB)

	//test
	// var whichUser string
	// var logInOrNot bool
	// rows, err = db.Query("SELECT username, loggedIn FROM users WHERE username = ?;", nicknameEmailDB)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	rows.Scan(&whichUser, &logInOrNot)
	// }
	// fmt.Printf("login user: %s, login status: %v\n", whichUser, logInOrNot)
}
