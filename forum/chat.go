package forum

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsChatResponse struct {
	Label       string   `json:"label"`
	Content     string   `json:"content"`
	OnlineUsers []string `json:"online_users"`
}

type WsChatPayload struct {
	Label      string `json:"label"`
	Content    string `json:"content"`
	ReceiverId int    `json:"receiver_id"` // for chat
}

var chatPayloadChan = make(chan WsChatPayload)

func chatWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Login Connected")
	var userListResponse WsChatResponse
	userListResponse.Label = "userList"

	rows, err := db.Query(`SELECT nickname FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var nicknameDBArr []string
	for rows.Next() {
		var nicknameDB string
		rows.Scan(&nicknameDB)
		nicknameDBArr = append(nicknameDBArr, nicknameDB)
	}
	fmt.Printf("nicknames: %v/n", nicknameDBArr)
	userListResponse.OnlineUsers = nicknameDBArr
	conn.WriteJSON(userListResponse)

	readChatPayloadFromWs(conn)
}

func readChatPayloadFromWs(conn *websocket.Conn) {
	defer func() {
		fmt.Println("Chat Ws Conn Closed")
	}()

	var chatPayload WsChatPayload
	for {
		err := conn.ReadJSON(&chatPayload)
		if err == nil {
			fmt.Printf("Sending chatPayload thru chan: %v\n", chatPayload)
			chatPayloadChan <- chatPayload
		}
	}
}
