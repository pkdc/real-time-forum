package forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type WsChatResponse struct {
	Label     string `json:"label"`
	Content   string `json:"content"`
	UserID    int    `json:"userID"`
	ContactID int    `json:"contactID"`
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
	// CookieValue string `json:"cookie_value"`
}

var (
	chatPayloadChan = make(chan WsChatPayload)
	ChatHub         *hub
	// chatWsMap       = make(map[*websocket.Conn]int)
	chatWsMap = make(map[int]*websocket.Conn)
)

func chatWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Chat Connected")

	// how to add client ws to chatWsMap?
	// is it easier to just use userlist ws for chat?
	// take up a space
	// chatWsMap[conn] = 0

	// create hub if none
	fmt.Printf("hub before %v\n", ChatHub)
	if ChatHub == nil { // if hub not made
		ChatHub = newHub()
		go ChatHub.Run()
	}
	fmt.Printf("hub after %v\n", ChatHub)

	// // find userID
	// find it
	// c, err := r.Cookie("session")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// var currentUserId int
	// rows, err := db.Query("SELECT userID, sessionID FROM sessions WHERE sessionID = ?;", c.Value)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	rows.Scan(&currentUserId)
	// }
	// // or get currentUserId by matching conn in userListWsMap
	// fmt.Printf("client conn in endpt: %v\n", conn)
	client := &Client{
		receiverRooms: make(map[int]*Room),
		// userID:        currentUserId,
		conn: conn,
		send: make(chan WsChatPayload),
	}
	// store conn in map?

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
	clientA  *Client
	clientB  *Client
	roomName string
}

var createRoomChan = make(chan roomRequest)

// create room when received roomRequest
func (h *hub) Run() {
	for {
		// readPump will send req thru createRoomChan if no existing room is found
		roomReq := <-createRoomChan
		fmt.Printf("the received room req: %v\n", roomReq)
		// var roomName string
		// // keep roomname in asc order
		// if participants.clientA.userID < participants.clientB.userID {
		// 	roomName = participants.clientA.userID) + "-and-" + participants.clientB.userID)
		// } else {
		// 	roomName = participants.clientB.userID) + "-and-" + participants.clientA.userID)
		// }
		rm := newRoom(roomReq.roomName, roomReq)
		fmt.Printf("created room name: %v\n", roomReq.roomName)
		go rm.run()
		// fmt.Printf("still fine\n")
		fmt.Printf(" hub: %v\n", h)
		fmt.Printf("new room in hub (before): %v\n", h.rooms[rm.roomName]) // deref ptr error
		h.rooms[roomReq.roomName] = rm
		fmt.Printf("new room in hub (after): %v\n", h.rooms[rm.roomName])
		// add room to reciverRooms (map) of clientA (c of c.readPump), feasible coz linked to c
		roomReq.clientA.receiverRooms[roomReq.clientB.userID] = rm
		fmt.Printf("%v has receiverRooms: %v\n", roomReq.clientA, roomReq.clientA.receiverRooms)
		// what if clientB initiate convo?
		// Now only clientA has this record, clientB doesn't have
		// prev checked if there is a room of this name? Add to clientB receiverRooms there?

		var createdRoomPayload WsChatPayload
		createdRoomPayload.Label = "created_room"
		createdRoomPayload.SenderId = roomReq.clientA.userID
		createdRoomPayload.ReceiverId = roomReq.clientB.userID
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
	fmt.Printf("room %v running\n", r)
	for {
		var chatRoomResponse WsChatPayload
		select {
		case chatRoomResponse = <-r.intoRoom:
			fmt.Printf("in room chatRoomResponse: %v\n", chatRoomResponse)
			// send to both clients when room receives msg
			fmt.Printf("sending chatRoomResponse %v to client: %v thru conn A %v\n", chatRoomResponse, r.clientA, r.clientA.conn)
			r.clientA.send <- chatRoomResponse
			fmt.Printf("sending chatRoomResponse %v to client: %v thru conn B %v\n", chatRoomResponse, r.clientB, r.clientB.conn)
			r.clientB.send <- chatRoomResponse
		}
	}
}

// func (r *Room) loadPrevMsgs() {
// }

// -----------------------Client-------------------------------
type Client struct {
	receiverRooms map[int]*Room
	userID        int
	conn          *websocket.Conn
	send          chan WsChatPayload // not added yet
}

// func idClient(conn *websocket.Conn) {
// 	// find userID
// 	c, err := r.Cookie("session")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	var currentUserId int
// 	rows, err := db.Query("SELECT userID, sessionID FROM sessions WHERE sessionID = ?;", c.Value)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		rows.Scan(&currentUserId)
// 	}
// 	// or get currentUserId by matching conn in userListWsMap

// 	client := &Client{
// 		receiverRooms: make(map[int]*Room),
// 		userID:        currentUserId,
// 		conn:          conn,
// 		send:          make(chan WsChatPayload),
// 	}
// 	fmt.Printf("Cleint created: %v\n", client)
// 	// go readChatPayloadFromWs(conn)
// 	go client.readPump()
// 	go client.writePump()
// }

func (c *Client) readPump() {
	defer func() {
		fmt.Println("readPump failed")
	}()

	fmt.Println("read pump running")
	var chatPayload WsChatPayload
	for {
		// fmt.Printf("client conn in readPump: %v\n", c.conn)
		err := c.conn.ReadJSON(&chatPayload)
		if err == nil {
			fmt.Printf("chatPayload %v with label:%s. of client %v \n", chatPayload, chatPayload.Label, c)
			fmt.Printf("chat err (should be nil): %v\n ", err)
			// create room
			if chatPayload.Label == "user-online" {
				// senderIdNum, err := chatPayload.SenderId)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				c.userID = chatPayload.SenderId
				fmt.Printf("chat: client has uid: %d\n", c.userID)
				chatWsMap[chatPayload.SenderId] = c.conn
				fmt.Printf("chat: added client ws to map, current map: %v\n", chatWsMap)
			} else if chatPayload.Label == "user-offline" {
				fmt.Printf("chat: user %d offline\n", loggedInUid)
				delete(chatWsMap, loggedInUid)
				fmt.Printf("chat: current map %v after logout\n", chatWsMap)
			} else if chatPayload.Label == "createChat" {
				// find the right room
				var findRoomName string
				if c.userID < chatPayload.ReceiverId {
					findRoomName = strconv.Itoa(c.userID) + "-and-" + strconv.Itoa(chatPayload.ReceiverId)
				} else {
					findRoomName = strconv.Itoa(chatPayload.ReceiverId) + "-and-" + strconv.Itoa(c.userID)
				}
				fmt.Printf("the right room name is: %s\n", findRoomName)
				rightChatRoom := ChatHub.findRoom(findRoomName) // can find right room
				fmt.Printf("the right room is: %v\n", rightChatRoom)

				if rightChatRoom == nil {
					fmt.Printf("right room is: not found \n")
					// if no record of the room, create one
					var rmReq roomRequest
					fmt.Printf("create rmReq \n")
					rmReq.roomName = findRoomName
					fmt.Printf("rmReq rn %s \n", rmReq.roomName)
					rmReq.clientA = c // link c and rmReq.clientA
					fmt.Printf("rmReq CA %v \n", rmReq.clientA)
					///////////////////////
					// dereference clientB IN THE RMReq, and put the userID or conn field into it
					// ReceiverIdNum, err := chatPayload.ReceiverId
					// if err != nil {
					// 	log.Fatal(err)
					// }

					clientB := &Client{
						receiverRooms: make(map[int]*Room),
						userID:        chatPayload.ReceiverId,
						conn:          chatWsMap[chatPayload.ReceiverId], // the receiver will also be online, so has a record in the map
						send:          make(chan WsChatPayload),
					}
					rmReq.clientB = clientB

					fmt.Printf("rmReq CB ID %v \n", rmReq.clientB.userID)
					fmt.Printf("rmReq CB conn %v \n", rmReq.clientB.conn)
					fmt.Printf("rmReq CB %v \n", rmReq.clientB)
					fmt.Printf("sending rmReq: %v\n", rmReq)
					createRoomChan <- rmReq
				}

				// // load the msg into rightChatRoom // not used yet
				// fmt.Println("----receiver", chatPayload.ReceiverId, "----sender", chatPayload.SenderId)
				// var creatingChatResponse WsChatResponse
				// // creatingChatResponse.Label= "using"
				// creatingChatResponse.Label = "chatBox"
				// // load prev msgs
				// // senderIdNum, _ := chatPayload.SenderId
				// // receiverIdNum, _ := chatPayload.ReceiverId
				// creatingChatResponse.Content = sortMessages(chatPayload.SenderId, chatPayload.ReceiverId)
				// // just loading for the sender!!
				// c.conn.WriteJSON(creatingChatResponse) // only writing to sender
				// // c.receiverRooms[chatPayload.ReceiverId]

				// reply? roomname?

			} else if chatPayload.Label == "chat" {
				fmt.Printf("Sending chatPayload thru chan: %v\n", chatPayload)
				var chatResponse WsChatPayload

				chatResponse.Label = "chat"
				chatResponse.Content = chatPayload.Content
				chatResponse.SenderId = chatPayload.SenderId
				chatResponse.ReceiverId = chatPayload.ReceiverId

				// if chatPayload.Online {
				// receiver online
				// send msg into room

				// finding the correct receiver room in the client
				// ReceiverIdNum, err := strconv.Atoi(chatPayload.ReceiverId)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				receivingRoom := *(c.receiverRooms[chatPayload.ReceiverId])
				fmt.Printf("chat sent to receivingRoom %v from client %v\n", receivingRoom, c)
				receivingRoom.intoRoom <- chatResponse
				// } else {
				// receiver offline
				// }
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
		chatResponse := <-c.send
		fmt.Printf("sneding payload %v to client %v\n", chatResponse, c)
		c.conn.WriteJSON(chatResponse)
	}
}

// ---------------------------------------
// func listeningChat(conn *websocket.Conn, msg WsChatPayload) {
// 	// var chatResponse WsChatResponse
// 	defer func() {
// 		fmt.Println("chat Ws Conn Closed")
// 	}()
// 	for {
// 		if msg.Label == "message" {
// 			var pureMsg WsChatPayload
// 			json.Unmarshal([]byte(msg.Content), &pureMsg)
// 			processMsg(pureMsg)
// 			fmt.Printf("payload received: %v\n", msg)
// 		}
// 	}
// }

// func processMsg(msg WsChatPayload) {
// 	rows, err := db.Prepare("INSERT INTO messages(senderID,receiverID,messageTime,content,seen) VALUES(?,?,?,?,?);")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	rows.Exec(msg.SenderId, msg.ReceiverId, time.Now(), msg.Content, false)
// 	fmt.Println("msg saved successfully")
// }
