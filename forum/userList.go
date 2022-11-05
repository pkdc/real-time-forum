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
	updateUList(conn)

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
			userListPayloadChan <- userListPayload
		}
	}
}

func ProcessAndReplyUserList() {
	receivedPayload := <-userListPayloadChan
	if receivedPayload.Label == "update" {
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
