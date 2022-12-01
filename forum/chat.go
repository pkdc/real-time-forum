package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type WsChatResponse struct {
	Label     string `json:"label"`
	Content   string `json:"content"`
	UserID    string `json:"userID"`
	ContactID string `json:"contactID"`
}
type MessageArray struct {
	Index int           `json:"index"`
	Msg   WsChatPayload `json:"msgInfo"`
}
type WsChatPayload struct {
	Label       string `json:"label"`
	Content     string `json:"content"`
	SenderId    int    `json:"sender_id"`
	ReceiverId  int    `json:"receiver_id"`
	Online      bool   `json:"online"` // whether the receiver is online
	MessageTime string `json:"message_time"`
	Noti        bool   `json:"noti"`
	Right       bool   `json:"right_side"`
}

var (
	chatPayloadChan = make(chan WsChatPayload)
	ChatHub         *hub
)

func chatWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Chat Connected")

	// create hub if none
	fmt.Printf("hub before %v\n", ChatHub)
	if ChatHub == nil { // if map not made
		ChatHub = newHub()
	}
	fmt.Printf("hub after %v\n", ChatHub)

	// find userID
	c, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		return
	}
	var currentUserId int
	rows, err := db.Query("SELECT userID, sessionID FROM sessions WHERE sessionID = ?;", c.Value)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&currentUserId)
	}
	// or get currentUserId by matching conn in userListWsMap

	client := &Client{
		receiverRooms: make(map[int]*Room),
		userID:        currentUserId,
		conn:          conn,
		send:          make(chan WsChatPayload),
	}
	fmt.Printf("Cleint created: %v\n", client)
	// go readChatPayloadFromWs(conn)
	go client.readPump()
	go client.writePump()
}

// go into readPump?
// func readChatPayloadFromWs(conn *websocket.Conn) {
// 	defer func() {
// 		fmt.Println("Chat Ws Conn Closed")
// 	}()

// 	var chatPayload WsChatPayload
// 	for {
// 		err := conn.ReadJSON(&chatPayload)
// 		if err == nil {

// 		}
// 	}
// }

// -----------------------Hub-------------------------------
// Hub to create rooms
// key roomname
type hub struct {
	rooms map[string]*Room // key: roomname
}

func newHub() *hub {
	return &hub{
		rooms: make(map[string]*Room),
	}
}

type roomRequest struct {
	clientA *Client
	clientB *Client
}

var createRoomChan = make(chan roomRequest)

// create room when received roomRequest
func (h *hub) Run() {
	for {
		participants := <-createRoomChan
		fmt.Printf("the received room req for: %v\n", participants)
		var roomName string
		// keep roomname in asc order
		if participants.clientA.userID < participants.clientB.userID {
			roomName = strconv.Itoa(participants.clientA.userID) + "-and-" + strconv.Itoa(participants.clientB.userID)
		} else {
			roomName = strconv.Itoa(participants.clientB.userID) + "-and-" + strconv.Itoa(participants.clientA.userID)
		}
		rm := newRoom(roomName, participants)
		fmt.Printf("created room name: %v\n", roomName)
		rm.run()
		h.rooms[roomName] = rm
		fmt.Printf("rooms in hub: %v\n", h.rooms)
		// add room to reciverRooms (map) of clientA (c of c.readPump), feasible coz linked to c
		participants.clientA.receiverRooms[participants.clientB.userID] = rm
		fmt.Printf("%v has receiverRooms: %v\n", participants.clientA, participants.clientA.receiverRooms)
		// what if clientB initiate convo?
		// Now only clientA has this record, clientB doesn't have
		// prev checked if there is a room of this name? Add to clientB receiverRooms there?

		var createdRoomPayload WsChatPayload
		createdRoomPayload.Label = "created_room"
		createdRoomPayload.SenderId = participants.clientA.userID
		createdRoomPayload.ReceiverId = participants.clientB.userID
		rm.intoRoom <- createdRoomPayload
	}
}

func (h *hub) findRoom(roomname string) *Room {
	elem, ok := ChatHub.rooms[roomname]
	if ok {
		return elem
	}
	return nil
}

// -----------------------Room-------------------------------
type Room struct {
	roomName string // eg: "1-and-2"
	clientA  *Client
	clientB  *Client
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
			// send to both clients when room receives msg
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

	fmt.Println("read pump running")
	var chatPayload WsChatPayload
	for {
		err := c.conn.ReadJSON(&chatPayload)
		if err == nil {
			fmt.Printf("reading payload %v from %v", chatPayload, c)
			// create room
			if chatPayload.Label == "room" {
				// find the right room
				var findRoomName string
				if chatPayload.SenderId < chatPayload.ReceiverId {
					findRoomName = strconv.Itoa(chatPayload.SenderId) + "-and-" + strconv.Itoa(chatPayload.ReceiverId)
				} else {
					findRoomName = strconv.Itoa(chatPayload.ReceiverId) + "-and-" + strconv.Itoa(chatPayload.SenderId)
				}
				rightChatRoom := ChatHub.findRoom(findRoomName)
				fmt.Printf("the right room is: %v", rightChatRoom)

				if rightChatRoom == nil {
					// if no record of the room, create one
					var rmReq roomRequest
					rmReq.clientA = c // link c and rmReq.clientA
					// dereference clientB and put the userID or conn field into it
					(*(rmReq.clientB)).userID = chatPayload.ReceiverId
					(*(rmReq.clientB)).conn = userListWsMap[chatPayload.ReceiverId]
					fmt.Printf("sending rmReq: %v\n", rmReq)
					createRoomChan <- rmReq
				}

				// load the msg into rightChatRoom
				fmt.Println("----receiver", chatPayload.ReceiverId, "----sender", chatPayload.SenderId)
				var creatingChatResponse WsChatResponse
				// creatingChatResponse.Label= "using"
				creatingChatResponse.Label = "prevMsgs"
				// load prev msgs
				creatingChatResponse.Content = sortMessages(chatPayload.SenderId, chatPayload.ReceiverId)
				// just loading for the sender!!
				c.conn.WriteJSON(creatingChatResponse)
				// c.receiverRooms[chatPayload.ReceiverId]

				// reply? roomname?

			} else if chatPayload.Label == "chat" {
				fmt.Printf("Sending chatPayload thru chan: %v\n", chatPayload)
				if chatPayload.Online {
					// receiver online
					// send msg into room
					// finding the correct room
					receivingRoom := *(c.receiverRooms[chatPayload.ReceiverId])
					receivingRoom.intoRoom <- chatPayload
				} else {
					// receiver offline
				}
			}
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		fmt.Println("writePump failed")
	}()
	fmt.Println("write pump running")
	for {
		chatPayload := <-c.send
		fmt.Printf("sneding payload %v to client", chatPayload)
		c.conn.WriteJSON(chatPayload)
	}
}

// ---------------------------------------
func listeningChat(conn *websocket.Conn, msg WsChatPayload) {
	// var chatResponse WsChatResponse
	defer func() {
		fmt.Println("chat Ws Conn Closed")
	}()
	for {
		if msg.Label == "message" {
			var pureMsg WsChatPayload
			json.Unmarshal([]byte(msg.Content), &pureMsg)
			processMsg(pureMsg)
			fmt.Printf("payload received: %v\n", msg)
		}
	}
}

func processMsg(msg WsChatPayload) {
	rows, err := db.Prepare("INSERT INTO messages(senderID,receiverID,messageTime,content,seen) VALUES(?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	rows.Exec(msg.SenderId, msg.ReceiverId, time.Now(), msg.Content, false)
	fmt.Println("msg saved successfully")
}
