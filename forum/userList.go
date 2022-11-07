package forum

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsUserListResponse struct {
	Label       string       `json:"label"`
	Content     string       `json:"content"`
	OnlineUsers []userStatus `json:"online_users"`
}

type WsUserListPayload struct {
	Label   string         `json:"label"`
	Content string         `json:"content"`
	Cookie  string         `json:"cookie"`
	Conn    websocket.Conn `json:"-"`
}

type userStatus struct {
	Nickname string `json:"nickname"`
	LoggedIn bool   `json:"status"`
}

var userListPayloadChan = make(chan WsUserListPayload)

func userListWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("UL Connected")

	readUserListPayloadFromWs(conn)
}

func readUserListPayloadFromWs(conn *websocket.Conn) {
	defer func() {
		fmt.Println("UserList Ws Conn Closed")
	}()

	var userListPayload WsUserListPayload
	for {
		err := conn.ReadJSON(&userListPayload)
		if err == nil {
			fmt.Printf("Sending userListPayload thru chan: %v\n", userListPayload)
			userListPayload.Conn = *conn
			userListPayloadChan <- userListPayload
		}
	}
}

func ProcessAndReplyUserList() {
	receivedUserListPayload := <-userListPayloadChan
	if receivedUserListPayload.Label == "update" {
		// can we get the cookie from backend directly?

		// // find which userID
		// c, err := r.Cookie("session")
		// if err != nil {
		// 	fmt.Println("User not logged in")
		// 	return
		// }
		// var loggedInUid int
		// rows, err := db.Query("SELECT userID FROM sessions WHERE sessionID = ?;", c.Value)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer rows.Close()
		// for rows.Next() {
		// 	rows.Scan(&loggedInUid)
		// }
		// fmt.Printf("loggedInUid UL %d \n", loggedInUid)
		// // store conn in websockets table
		// stmt, err := db.Prepare(`INSERT INTO websockets
		// 					(userID, websocketAdd, usage)
		// 					VALUES (?, ?, ?);`)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer stmt.Close()
		// stmt.Exec(loggedInUid, receivedUserListPayload.Conn, "userlist")
		updateUList()
	}
}

func updateUList() {
	var userListResponse WsUserListResponse
	userListResponse.Label = "reg"

	rows, err := db.Query(`SELECT nickname, loggedIn FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var userStatusDBArr []userStatus
	for rows.Next() {
		var nicknameDB string
		var loggedInDB bool
		rows.Scan(&nicknameDB, &loggedInDB)
		userStatusElement := struct {
			Nickname string `json:"nickname"`
			LoggedIn bool   `json:"status"`
		}{
			nicknameDB,
			loggedInDB,
		}
		userStatusDBArr = append(userStatusDBArr, userStatusElement)
	}

	fmt.Printf("nicknames: %v\n", userStatusDBArr)
	userListResponse.OnlineUsers = userStatusDBArr
	broadcast(userListResponse)
	// conn.WriteJSON(userListResponse)
}

func broadcast(userListResponse WsUserListResponse) {
	rows, err := db.Query(`SELECT websocketAdd FROM websockets`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var wsArr []*websocket.Conn
	for rows.Next() {
		var wsAdd *websocket.Conn
		rows.Scan(&wsAdd)
		wsArr = append(wsArr, wsAdd)
	}
	for _, wsAddress := range wsArr {
		wsAddress.WriteJSON(userListResponse)
	}
}
