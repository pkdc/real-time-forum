package forum

import (
	"fmt"
	"log"
	"net/http"
)

type WsUserListResponse struct {
	Label       string       `json:"label"`
	Content     string       `json:"content"`
	OnlineUsers []userStatus `json:"online_users"`
}

type WsUserListPayload struct {
	Label   string `json:"label"`
	Content string `json:"content"`
}

type userStatus struct {
	Nickname string `json:"nickname"`
	LoggedIn bool   `json:"status"`
}

// var chatPayloadChan = make(chan WsChatPayload)

func userListWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("UL Connected")
	var userListResponse WsUserListResponse
	userListResponse.Label = "initial"

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
	conn.WriteJSON(userListResponse)

	// readChatPayloadFromWs(conn)
}

// func readChatPayloadFromWs(conn *websocket.Conn) {
// 	defer func() {
// 		fmt.Println("Chat Ws Conn Closed")
// 	}()

// 	var chatPayload WsChatPayload
// 	for {
// 		err := conn.ReadJSON(&chatPayload)
// 		if err == nil {
// 			fmt.Printf("Sending chatPayload thru chan: %v\n", chatPayload)
// 			chatPayloadChan <- chatPayload
// 		}
// 	}
// }
