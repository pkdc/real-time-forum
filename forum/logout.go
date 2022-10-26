package forum

import (
	"fmt"
	"log"
	"net/http"
)

type WsLogoutResponse struct {
	Label   string `json:"label"`
	Content string `json:"content"`
}

// func LogoutWsEndpoint(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	fmt.Println("Logout WS Connected")
// 	var firstResponse WsLogoutResponse
// 	firstResponse.Label = "logout"
// 	firstResponse.Content = "Thanks for visiting"
// 	conn.WriteJSON(firstResponse)

// 	removeCookie()
// }

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	var logoutUname string

	if err == nil {
		// get the username of the logout user
		rows, err := db.Query("SELECT username FROM sessions WHERE sessionID = ?;", c.Value)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&logoutUname)
		}
		fmt.Printf("Found user %s wants to logout", logoutUname)

		// delete sessionID from sessions db table
		stmt, err := db.Prepare("DELETE FROM sessions WHERE sessionID=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(c.Value)
		fmt.Printf("cookie sid removed (have value): %s\n", c.Value)
	}

	// test
	var sessionID string
	rows, err := db.Query("SELECT * FROM sessions")
	for rows.Next() {
		rows.Scan(&sessionID)
	}
	fmt.Printf("cookie sid removed (should be empty): %s\n", sessionID) // empty is correct

	// delete browser's cookie
	_, err = r.Cookie("session")
	if err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		})
	}
	fmt.Printf("%s Logout\n", logoutUname)

	stmt, err := db.Prepare("UPDATE users SET loggedIn = ? WHERE username = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(false, logoutUname)
}
