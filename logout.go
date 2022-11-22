package main

import (
	"fmt"
	"log"
	"net/http"
)

func processLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout handler reached")
	c, err := r.Cookie("session")
	var logoutUid int

	if err == nil {
		// get the uid of the logout user
		rows, err := db.Query("SELECT userID FROM sessions WHERE sessionID = ?;", c.Value)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&logoutUid)
		}
		fmt.Printf("Found userID %d wants to logout\n", logoutUid)

		// delete sessionID from sessions db table
		// still need sessionID record for removing logout user from user list
		// stmt, err := db.Prepare("DELETE FROM sessions WHERE sessionID=?")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer stmt.Close()
		// stmt.Exec(c.Value)
		// fmt.Printf("cookie sid removed (have value): %s\n", c.Value)
	}

	// // test
	// var sessionID string
	// rows, err := db.Query("SELECT * FROM sessions")
	// for rows.Next() {
	// 	rows.Scan(&sessionID)
	// }
	// fmt.Printf("cookie sid removed (should be empty): %s\n", sessionID) // empty is correct

	fmt.Printf("%d Logout\n", logoutUid)

	stmt, err := db.Prepare("UPDATE users SET loggedIn = ? WHERE userID = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(false, logoutUid)
}
