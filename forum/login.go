package forum

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsLoginResponse struct {
	Label   string `json:"label"`
	Content string `json:"content"`
}
type WsLoginPayload struct {
	Label    string          `json:"label"`
	Name     string          `json:"name"`
	Password string          `json:"pw"`
	Conn     *websocket.Conn `json:"-"`
}

func LoginWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Connected")
	var response WsLoginResponse
	response.Label = "Greet"
	response.Content = "Welcome to the Forum!"
	conn.WriteJSON(response)
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
			loginPayload.Conn = conn
			fmt.Printf("payload: %v\n", loginPayload)
			loginPayloadChan <- loginPayload
		}
	}
}

func ListenToLoginChan() {
	for {
		e := <-loginPayloadChan
		fmt.Println("payload received: ", e)
	}

}
