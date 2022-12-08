package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	LoadMsg     bool           `json:"loadMsg"`
}

type userStatus struct {
	Nickname   string `json:"nickname"`
	LoggedIn   bool   `json:"status"`
	UserID     int    `json:"userID"`
	MsgCheck   bool   `json:"msgcheck"`
	CurUser    bool   `json:"curuser"`
	withoutlet bool
}

var (
	PageMsgMap          = make(map[int]int)
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
		// fmt.Println("Label", userListPayload.Label)
		if err == nil && userListPayload.Label == "createChat" {
			fmt.Println("----contact", userListPayload.ContactID, "----userID", userListPayload.UserID)
			var creatingChatResponse WsUserListResponse
			// creatingChatResponse.Label= "using"
			creatingChatResponse.Label = "chatBox"
			if !userListPayload.LoadMsg {
				PageMsgMap[userListPayload.UserID] = 1234567890
			}
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
		PageMsgMap[loggedInUid] = 1234567890
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
			// store conn in map
			userListWsMap[loggedInUid] = &receivedUserListPayload.Conn
			fmt.Printf("UL current map: %v", userListWsMap)
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
	var tempUserStatus []userStatus
	var userStatusDBArr []userStatus
	for rows.Next() {
		var nicknameDB string
		var loggedInDB bool
		var UserIDDB int
		var msgcheck bool
		rows.Scan(&nicknameDB, &loggedInDB, &UserIDDB)
		userStatusElement := struct {
			Nickname   string `json:"nickname"`
			LoggedIn   bool   `json:"status"`
			UserID     int    `json:"userID"`
			MsgCheck   bool   `json:"msgcheck"`
			CurUser    bool   `json:"curuser"`
			withoutlet bool
		}{
			nicknameDB,
			loggedInDB,
			UserIDDB,
			msgcheck,
			false,
			false,
		}
		tempUserStatus = append(tempUserStatus, userStatusElement)
	}
	topOfTheList := sortConversations()

	var letter []userStatus
	var notLetter []userStatus
	var msgHistory []userStatus
	for i := 0; i < len(tempUserStatus); i++ {
		for k := 0; k < len(topOfTheList); k++ {
			if tempUserStatus[i].UserID == loggedInUid {
				tempUserStatus[i].CurUser = true
				continue
			}
			if tempUserStatus[i].UserID == topOfTheList[k] {
				tempUserStatus[i].MsgCheck = true
				tempUserStatus[k], tempUserStatus[i] = tempUserStatus[i], tempUserStatus[k]
			}
		}
	}
	for i := 0; i < len(tempUserStatus); i++ {
		if strings.Title(tempUserStatus[i].Nickname)[0] < 64 || strings.Title(tempUserStatus[i].Nickname)[0] > 91 {
			tempUserStatus[i].withoutlet = true
		}
		if tempUserStatus[i].MsgCheck {
			msgHistory = append(msgHistory, tempUserStatus[i])
		} else if !tempUserStatus[i].withoutlet {
			letter = append(letter, tempUserStatus[i])
		} else {
			notLetter = append(notLetter, tempUserStatus[i])
		}

	}
	counter := 0
loop:
	for i := 0; i < len(letter)-1; i++ {
		if strings.Title(letter[i].Nickname)[0] > strings.Title(letter[i+1].Nickname)[0] {
			letter[i], letter[i+1] = letter[i+1], letter[i]
		}
		if counter != 2 && i == len(letter)-2 {
			counter++
			goto loop
		}
	}
	counter2 := 0
loop2:
	for i := 0; i < len(notLetter)-1; i++ {
		if strings.Title(notLetter[i].Nickname)[0] > strings.Title(notLetter[i+1].Nickname)[0] {
			notLetter[i], notLetter[i+1] = notLetter[i+1], notLetter[i]
		}
		if counter != 2 && i == len(notLetter)-2 {
			counter2++
			goto loop2
		}
	}
	userStatusDBArr = append(userStatusDBArr, msgHistory...)
	userStatusDBArr = append(userStatusDBArr, letter...)
	userStatusDBArr = append(userStatusDBArr, notLetter...)

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
	fmt.Println(PageMsgMap)
	rows, err := db.Query(
		`SELECT * 
	FROM messages 
	WHERE messageID < ? AND ((senderID = ? AND receiverID = ?) OR (receiverID = ? AND senderID = ?))
	ORDER BY messageID DESC	
	LIMIT ?
	;`, PageMsgMap[sendID], sendID, recID, sendID, recID, 10)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var oneMsg WsChatPayload
		var msgTime time.Time
		var msgID int
		rows.Scan(&msgID, &(oneMsg.SenderId), &(oneMsg.ReceiverId), &msgTime, &(oneMsg.Content), &(oneMsg.Noti))
		oneMsg.MessageTime = msgTime.String()
		if oneMsg.SenderId == loggedInUid {
			oneMsg.Right = true
		}
		fmt.Println(msgID, "INDEX")
		allMsg.Index = msgID
		allMsg.Msg = oneMsg
		arrMsgArray = append(arrMsgArray, allMsg)
		PageMsgMap[sendID] = msgID
	}
	fmt.Println("chatinfo:", arrMsgArray)

	return arrMsgArray
}

func sortMessages(sendID, recID int) string {
	// firstMes := displayChatInfo(sendID, recID)
	// secMes := displayChatInfo(recID, sendID)
	// allMes := append(firstMes, secMes...)

	allMes := displayChatInfo(sendID, recID)
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

func sortConversations() []int {
	var allCon []int
	rows, err := db.Query("SELECT receiverID FROM messages WHERE senderID= ?;", loggedInUid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var recID int
		rows.Scan(&recID)
		allCon = append(allCon, recID)

	}
	for i := 0; i < len(allCon)/2; i++ {
		allCon[i], allCon[len(allCon)-(i+1)] = allCon[len(allCon)-(i+1)], allCon[i]
	}
	var lastOne []int
	for _, v := range allCon {
		skip := false
		for _, u := range lastOne {
			if v == u {
				skip = true
				break
			}
		}
		if !skip {
			lastOne = append(lastOne, v)
		}
	}
	return lastOne
}
