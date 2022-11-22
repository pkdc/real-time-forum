package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type WsChatResponse struct {
	Label   string `json:"label"`
	Content string `json:"content"`
}

type WsChatPayload struct {
	Label      string `json:"label"`
	Content    string `json:"content"`
	SenderId   int    `json:"sender_id"`
	ReceiverId int    `json:"receiver_id"`
	Online     bool   `json:"online"` // whether the receiver is online
}

var chatPayloadChan = make(chan WsChatPayload)

func chatWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Chat Connected %v", conn)

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
			// find the right room
			var findRoomName string
			if chatPayload.SenderId < chatPayload.ReceiverId {
				findRoomName = strconv.Itoa(chatPayload.SenderId) + "-and-" + strconv.Itoa(chatPayload.ReceiverId)
			} else {
				findRoomName = strconv.Itoa(chatPayload.ReceiverId) + "-and-" + strconv.Itoa(chatPayload.SenderId)
			}
			hub.findRoom(findRoomName)

			// load the msg

			// if no record of the room
			var rmReq roomRequest
			rmReq.clientA.userID = chatPayload.SenderId
			rmReq.clientB.userID = chatPayload.ReceiverId
			userListWsMap[rmReq.clientA.userID] = rmReq.clientA.conn
			userListWsMap[rmReq.clientB.userID] = rmReq.clientB.conn
			createRoomChan <- rmReq

			fmt.Printf("Sending chatPayload thru chan: %v\n", chatPayload)
			if chatPayload.Online {
				//  receiver online
				// if chatPayload.Label == "create-room" {

				// }
			} else {
				// receiver offline
			}
			// listeningChat(conn, chatPayload)
			// chatPayloadChan <- chatPayload
		}
	}
}

// -----------------------Hub-------------------------------
// Hub to create rooms
// key roomname
type Hub struct {
	rooms map[string]Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]Room),
	}
}

type roomRequest struct {
	clientA Client
	clientB Client
}

var createRoomChan = make(chan roomRequest)

// create room when received roomRequest
func (h *Hub) Run() {
	for {
		participants := <-createRoomChan
		var roomName string
		// keep roomname in asc order
		if participants.clientA.userID < participants.clientB.userID {
			roomName = strconv.Itoa(participants.clientA.userID) + "-and-" + strconv.Itoa(participants.clientB.userID)
		} else {
			roomName = strconv.Itoa(participants.clientB.userID) + "-and-" + strconv.Itoa(participants.clientA.userID)
		}
		rm := newRoom(roomName, participants)
		h.rooms[roomName] = *rm
		// h.rooms = append(h.rooms, *rm)
	}
}

func (h *Hub) findRoom(roomname string) {

}

// -----------------------Room-------------------------------
type Room struct {
	roomName string // eg: "1-and-2"
	clientA  Client
	clientB  Client
	intoRoom chan WsChatPayload
}

func newRoom(roomName string, participants roomRequest) *Room {
	return &Room{
		roomName: roomName,
		clientA:  participants.clientA,
		clientB:  participants.clientB,
		intoRoom: make(chan WsChatPayload),
	}
}

// -----------------------Client-------------------------------
type Client struct {
	userID int
	conn   *websocket.Conn
	send   chan WsChatPayload // not added yet
}

func (c *Client) readPump() {
	var chatPayload WsChatPayload
	for {
		err := c.conn.ReadJSON(&chatPayload)
		if err == nil {
			// how to find out which room?
			// intoRoom <- chatPayload
		}
	}

}

func (c *Client) writePump() {

}
