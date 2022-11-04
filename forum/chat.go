package forum

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsChatResponse struct {
	Label   string `json:"label"`
	Content string `json:"content"`
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
	fmt.Println("Chat Connected")

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
