package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type WsLoginResponse struct {
	Label   string `json:"label"`
	Content string `json:"content"`
	Pass    bool   `json:"pass"`
}
type WsLoginPayload struct {
	Label         string `json:"label"`
	NicknameEmail string `json:"name"`
	Password      string `json:"pw"`
	// Conn          *websocket.Conn `json:"-"`
}

func LoginWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Connected")
	var firstResponse WsLoginResponse
	firstResponse.Label = "Greet"
	firstResponse.Content = "Please login to the Forum"
	conn.WriteJSON(firstResponse)
	// insert conn into db with empty userID, fill in the userID when registered or logged in
	// stmt, err := db.Prepare(`INSERT INTO websockets (userID, websocketAdd) VALUES (?, ?);`)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// stmt.Exec("", conn)
	listenToLoginWs(conn)
}

func listenToLoginWs(conn *websocket.Conn) {
	defer func() {
		fmt.Println("Ws Conn Closed")
	}()

	var loginPayload WsLoginPayload

	for {
		err := conn.ReadJSON(&loginPayload)
		if err == nil {
			// loginPayload.Conn = conn
			fmt.Printf("payload: %v\n", loginPayload)
		}
		ProcessAndReplyLogin(conn, loginPayload)
	}
}

func ProcessAndReplyLogin(conn *websocket.Conn, loginPayload WsLoginPayload) {

	fmt.Println("payload received: ", loginPayload)

	if loginPayload.Label == "login" {
		fmt.Printf("login u: %s: , login pw: %s\n", loginPayload.NicknameEmail, loginPayload.Password)

		// // get user data from db
		var nicknameEmailDB string
		var hashDB []byte

		fmt.Printf("%s trying to Login\n", loginPayload.NicknameEmail)
		rows, err := db.Query(`SELECT nickname, email, password 
							FROM users
							WHERE nickname = ?
							OR email = ?`, loginPayload.NicknameEmail, loginPayload.NicknameEmail)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&nicknameEmailDB, &hashDB)
		}

		// // test hash
		// hash, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
		// fmt.Printf("nicknameEmailDB: %s , hashDB: %s\n", nicknameEmailDB, hashDB)

		// // compare pw
		err = bcrypt.CompareHashAndPassword(hashDB, []byte(loginPayload.Password))
		// fmt.Printf("DB pw: %s, entered: %s\n", hashDB, loginPayload.password)
		// fmt.Printf("DB pw: %s, entered: %s\n", hashDB, hash)

		// Login failed
		if err != nil {
			// login failed
			var failedResponse WsLoginResponse
			failedResponse.Label = "login"
			failedResponse.Pass = false
			resp, err := json.Marshal(failedResponse)
			if err != nil {
				log.Fatal(err)
			}
			conn.WriteJSON(resp)
			return
		}
		// Login successfully
		fmt.Printf("%s (name from DB) Login successfully\n", loginPayload.NicknameEmail)

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

		// // assign a cookie
		// sid := uuid.NewV4()
		// fmt.Printf("login sid: %s\n", sid)
		// http.SetCookie(w, &http.Cookie{
		// 	Name:   "session",
		// 	Value:  sid.String(),
		// 	MaxAge: 900, // 15mins
		// })

		// // forumUser.Username = nicknameEmailDB
		// // forumUser.LoggedIn = true
		// // fmt.Printf("%s forum User Login\n", forumUser.Username)

		// // update the user's login status
		stmt, err := db.Prepare("UPDATE users SET loggedIn = ? WHERE nickname = ? OR email = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(true, nicknameEmailDB, nicknameEmailDB)

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

}
