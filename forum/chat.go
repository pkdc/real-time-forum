package forum

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
var ChatHub *hub

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

	fmt.Printf("hub before %v", ChatHub)
	if (*ChatHub).rooms == nil { // if map not made
		ChatHub = newHub()
	}
	fmt.Printf("hub after %v", ChatHub)

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
			rightChatRoom := ChatHub.findRoom(findRoomName)
			if rightChatRoom == nil {
				// if no record of the room
				var rmReq roomRequest
				rmReq.clientA.userID = chatPayload.SenderId
				rmReq.clientB.userID = chatPayload.ReceiverId
				userListWsMap[rmReq.clientA.userID] = rmReq.clientA.conn
				userListWsMap[rmReq.clientB.userID] = rmReq.clientB.conn
				createRoomChan <- rmReq
			} else {
				// load the msg into rightChatRoom
			}

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
type hub struct {
	rooms map[string]Room
}

func newHub() *hub {
	return &hub{
		rooms: make(map[string]Room),
	}
}

type roomRequest struct {
	clientA Client
	clientB Client
}

var createRoomChan = make(chan roomRequest)

// create room when received roomRequest
func (h *hub) Run() {
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

func (h *hub) findRoom(roomname string) *Room {
	elem, ok := ChatHub.rooms[roomname]
	if ok {
		return &elem
	}
	return nil
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

func (r *Room) run() {
	for {
		var chatRoomPayload WsChatPayload
		select {
		case chatRoomPayload = <-r.intoRoom:
			fmt.Printf("in room chatRoomPayload: %v", chatRoomPayload)
			// send to both
			r.clientA.send <- chatRoomPayload
			r.clientB.send <- chatRoomPayload
		}
	}
}

// -----------------------Client-------------------------------
type Client struct {
	receiverRooms map[int]*Room
	userID        int
	conn          *websocket.Conn
	send          chan WsChatPayload // not added yet
}

func (c *Client) readPump() {
	defer func() {
		fmt.Println("readPump failed")
	}()

	var chatPayload WsChatPayload
	for {
		err := c.conn.ReadJSON(&chatPayload)
		if err == nil {
			// find out which room
			receivingRoom := *(c.receiverRooms[chatPayload.ReceiverId])
			receivingRoom.intoRoom <- chatPayload

		}
	}

}

func (c *Client) writePump() {
	defer func() {
		fmt.Println("writePump failed")
	}()

}
