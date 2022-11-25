package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WsUserListResponse struct {
	Label       string       `json:"label"`
	Content     string       `json:"content"`
	OnlineUsers []userStatus `json:"online_users"`
}

type WsUserListPayload struct {
	Label       string         `json:"label"`
	Content     string         `json:"content"`
	CookieValue string         `json:"cookie_value"`
	Conn        websocket.Conn `json:"-"`
	ContactID   int            `json:"contactID"`
	UserID      int            `json:"userID"`
}

type userStatus struct {
	Nickname string `json:"nickname"`
	LoggedIn bool   `json:"status"`
	UserID   int    `json:"userID"`
}

var (
	userListPayloadChan = make(chan WsUserListPayload)
	userListWsMap       = make(map[int]*websocket.Conn)
	loggedInUid         int
)

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
		// fmt.Print("ul ")
		err := conn.ReadJSON(&userListPayload)
		fmt.Println("Label",userListPayload.Label)
		fmt.Println(err)
		if err == nil && userListPayload.Label == "createChat" {
			fmt.Println("----contact", userListPayload.ContactID, "----userID", userListPayload.UserID)
			var creatingChatResponse WsUserListResponse
			// creatingChatResponse.Label= "using"
			creatingChatResponse.Label = "chatBox"
			creatingChatResponse.Content = sortMessages(userListPayload.UserID, userListPayload.ContactID)
			conn.WriteJSON(creatingChatResponse)
		} else if err == nil {
			fmt.Printf("Sending userListPayload thru chan: %v\n", userListPayload)
			userListPayload.Conn = *conn
			userListPayloadChan <- userListPayload
		}
	}
}

func ProcessAndReplyUserList() {
	for {
		receivedUserListPayload := <-userListPayloadChan
		// payloadLabels := strings.Split(receivedUserListPayload.Label, "-")

		// find which userID
		rows, err := db.Query("SELECT userID FROM sessions WHERE sessionID = ?;", receivedUserListPayload.CookieValue)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&loggedInUid)
		}
		fmt.Printf("loggedInUid UL %d \n", loggedInUid)

		// close and remove conn from map if logout
		// if len(payloadLabels) > 1 && payloadLabels[1] == "logout" {
		if receivedUserListPayload.Label == "logout-update" {
			// _ = receivedUserListPayload.Conn.Close()
			delete(userListWsMap, loggedInUid)
			fmt.Print("----------------------------\n")
			fmt.Printf("removing  logout user: %d\n", loggedInUid)
			// fmt.Printf("removing ws logout user: %v\n", &receivedUserListPayload.Conn)
			fmt.Printf("uList after removing logout user: %v\n", userListWsMap)

			// delete sessionID from sessions db table
			// still need sessionID record for removing logout user from user list
			stmt, err := db.Prepare("DELETE FROM sessions WHERE sessionID=?")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			stmt.Exec(receivedUserListPayload.CookieValue)
			fmt.Printf("cookie sid removed (have value): %s\n", receivedUserListPayload.CookieValue)
		}

		// if len(payloadLabels) == 1 && payloadLabels[0] == "update" {
		if receivedUserListPayload.Label == "login-reg-update" {
			// store conn in websockets table
			// stmt, err := db.Prepare(`INSERT INTO websockets
			// 					(userID, websocketAdd, usage)
			// 					VALUES (?, ?, ?);`)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// defer stmt.Close()
			// fmt.Printf("uid: %d, Conn: %v, usage %s \n", loggedInUid, receivedUserListPayload.Conn, "userlist")
			// stmt.Exec(loggedInUid, receivedUserListPayload.Conn, "userlist")

			// store conn in map
			// userListWsMap[&receivedUserListPayload.Conn] = loggedInUid
			userListWsMap[loggedInUid] = &receivedUserListPayload.Conn
			fmt.Printf("current map: %v", userListWsMap)
		}
		updateUList()
	}
}

func updateUList() {
	var userListResponse WsUserListResponse
	userListResponse.Label = "update"

	rows, err := db.Query(`SELECT nickname, loggedIn, userID   FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var userStatusDBArr []userStatus
	for rows.Next() {
		var nicknameDB string
		var loggedInDB bool
		var UserIDDB int

		rows.Scan(&nicknameDB, &loggedInDB, &UserIDDB)
		fmt.Println(UserIDDB)
		userStatusElement := struct {
			Nickname string `json:"nickname"`
			LoggedIn bool   `json:"status"`
			UserID   int    `json:"userID"`
		}{
			nicknameDB,
			loggedInDB,
			UserIDDB,
		}
		userStatusDBArr = append(userStatusDBArr, userStatusElement)
	}

	fmt.Printf("UL nicknames: %v\n", userStatusDBArr)
	userListResponse.OnlineUsers = userStatusDBArr
	broadcast(userListResponse)
}

func broadcast(userListResponse WsUserListResponse) {
	for _, userListWs := range userListWsMap {
		userListWs.WriteJSON(userListResponse)
	}
	// rows, err := db.Query(`SELECT websocketAdd FROM websockets`) // WHERE usage = userlist
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// var wsArr []*websocket.Conn
	// for rows.Next() {
	// 	var wsAdd *websocket.Conn
	// 	rows.Scan(&wsAdd)
	// 	wsArr = append(wsArr, wsAdd)
	// }
	// for _, wsAddress := range wsArr {
	// 	wsAddress.WriteJSON(userListResponse)
	// }
}

func displayChatInfo(sendID, recID int) []MessageArray {
	var allMsg MessageArray
	var arrMsgArray []MessageArray
	rows, err := db.Query(
		`SELECT * 
	FROM messages 
	WHERE senderID = ? 
	AND receiverID = ?`, sendID, recID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var oneMsg WsChatPayload
		var msgTime time.Time
		var msgID int
		rows.Scan(&msgID, &(oneMsg.SenderId), &(oneMsg.ReceiverId), &msgTime, &(oneMsg.Content), &(oneMsg.Noti))
		fmt.Println("dont be empty", oneMsg.Content, len(oneMsg.Content))
		oneMsg.MessageTime = msgTime.String()
		fmt.Println(oneMsg.SenderId, "-----", loggedInUid)
		if oneMsg.SenderId == loggedInUid {
			oneMsg.Right = true
		}
		allMsg.Index = msgID
		allMsg.Msg = oneMsg
		arrMsgArray = append(arrMsgArray, allMsg)
	}
	fmt.Println("chatinfo:", arrMsgArray)

	return arrMsgArray
}

func sortMessages(sendID, recID int) string {
	firstMes := displayChatInfo(sendID, recID)
	secMes := displayChatInfo(recID, sendID)
	allMes := append(firstMes, secMes...)
	for k := 0; k < 10; k++ {
		for i := 0; i < len(allMes)-1; i++ {
			if allMes[i].Index > allMes[i+1].Index {
				allMes[i], allMes[i+1] = allMes[i+1], allMes[i]
			}
		}
	}
	jsonF, err := json.Marshal(allMes)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonF)
}
