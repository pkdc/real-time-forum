package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type WsUserListResponse struct {
	Label       string       `json:"label"`
	Content     string       `json:"content"`
	RealUser    int          `json:"realUser"`
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
	Nickname     string `json:"nickname"`
	LoggedIn     bool   `json:"status"`
	UserID       int    `json:"userID"`
	MsgCheck     bool   `json:"msgcheck"`
	Notification string `json:"noti"`
	withoutlet   bool
}

type userSort struct {
	usID  int
	msgID int
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
			var creatingChatResponse WsUserListResponse
			// creatingChatResponse.Label= "using"
			creatingChatResponse.Label = "chatBox"
			if !userListPayload.LoadMsg {
				PageMsgMap[userListPayload.UserID] = 1234567890
			}
			ChangeNotif(userListPayload.UserID, userListPayload.ContactID)
			creatingChatResponse.Content = sortMessages(userListPayload.UserID, userListPayload.ContactID)
			creatingChatResponse.RealUser = userListPayload.UserID
			conn.WriteJSON(creatingChatResponse)
			updateUList()
		} else if err == nil {
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
			fmt.Println("continue")
		}
		updateUList()
	}
}

func updateUList() {
	var userListResponse WsUserListResponse
	userListResponse.Label = "update"

	rows, err := db.Query(`SELECT nickname, loggedIn, userID ,notifications  FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var tempUserStatus []userStatus
	for rows.Next() {
		var nicknameDB string
		var loggedInDB bool
		var UserIDDB int
		var msgcheck bool
		var notifications string
		rows.Scan(&nicknameDB, &loggedInDB, &UserIDDB, &notifications)
		userStatusElement := struct {
			Nickname     string `json:"nickname"`
			LoggedIn     bool   `json:"status"`
			UserID       int    `json:"userID"`
			MsgCheck     bool   `json:"msgcheck"`
			Notification string `json:"noti"`
			withoutlet   bool
		}{
			nicknameDB,
			loggedInDB,
			UserIDDB,
			msgcheck,
			notifications,
			false,
		}
		tempUserStatus = append(tempUserStatus, userStatusElement)
	}
	for ind, x := range tempUserStatus {
		fmt.Println("index:", ind, "user:", x)
	}
	userListResponse.OnlineUsers = UserListSort(tempUserStatus)
	userListResponse.RealUser = loggedInUid
	broadcast(userListResponse)
}

func UserListSort(tempUserStatus []userStatus) []userStatus {
	var userStatusDBArr []userStatus
	topOfTheList := sortConversations()
	var letter []userStatus
	var notLetter []userStatus
	var msgHistory []userStatus
	var counter int
	var msgcheckcounter int
	for i := 0; i < len(tempUserStatus); i++ {
		for k := 0; k < len(topOfTheList); k++ {
			if tempUserStatus[i].UserID == topOfTheList[k] {
				tempUserStatus[i].MsgCheck = true
				// tempUserStatus[k], tempUserStatus[i] = tempUserStatus[i], tempUserStatus[k]
			}
		}
	}
	for ind, x := range tempUserStatus {
		fmt.Println("index:", ind, "user:", x.MsgCheck)
	}
	for i := 0; i < len(tempUserStatus); i++ {
		if strings.Title(tempUserStatus[i].Nickname)[0] < 64 || strings.Title(tempUserStatus[i].Nickname)[0] > 91 {
			tempUserStatus[i].withoutlet = true
		}
		if tempUserStatus[i].MsgCheck {
			msgHistory = append(msgHistory, tempUserStatus[topOfTheList[msgcheckcounter]-1])
			msgcheckcounter++
			fmt.Println(msgHistory, msgcheckcounter, "---msgHistor")
		} else if !tempUserStatus[i].withoutlet {
			letter = append(letter, tempUserStatus[i])
		} else {
			notLetter = append(notLetter, tempUserStatus[i])
		}

	}

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
		if counter2 != 2 && i == len(notLetter)-2 {
			counter2++
			goto loop2
		}
	}
	userStatusDBArr = append(userStatusDBArr, msgHistory...)
	userStatusDBArr = append(userStatusDBArr, letter...)
	userStatusDBArr = append(userStatusDBArr, notLetter...)
	fmt.Println("last status", userStatusDBArr)
	// fmt.Printf("UL nicknames: %v\n", userStatusDBArr)
	return userStatusDBArr
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
		var msgID int
		rows.Scan(&msgID, &(oneMsg.SenderId), &(oneMsg.ReceiverId), &(oneMsg.MessageTime), &(oneMsg.Content), &(oneMsg.Noti))
		// rows2, err := db.Prepare("UPDATE users SET seen = ? WHERE msgID = ?;")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer rows2.Close()
		// rows2.Exec(true, msgID)
		if oneMsg.SenderId == sendID {
			oneMsg.Right = true
		}
		allMsg.Index = msgID
		allMsg.Msg = oneMsg
		arrMsgArray = append(arrMsgArray, allMsg)
		PageMsgMap[sendID] = msgID
	}
	return arrMsgArray
}

func sortMessages(sendID, recID int) string {
	allMes := displayChatInfo(sendID, recID)
	jsonF, err := json.Marshal(allMes)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonF)
}

func sortConversations() []int {
	var allCon []int
	var allCon3 []userSort
	rows, err := db.Query("SELECT receiverID,messageID FROM messages WHERE senderID= ?;", loggedInUid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var allCon2 userSort
		rows.Scan(&allCon2.usID, &allCon2.msgID)
		allCon3 = append(allCon3, allCon2)

	}
	rows2, err2 := db.Query("SELECT senderID,messageID FROM messages WHERE receiverID= ?;", loggedInUid)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer rows2.Close()
	for rows2.Next() {
		var allCon2 userSort
		rows2.Scan(&allCon2.usID, &allCon2.msgID)
		allCon3 = append(allCon3, allCon2)

	}
	for k := 0; k < 100; k++ {
		for i := 0; i < len(allCon3)-1; i++ {
			if allCon3[i].msgID < allCon3[i+1].msgID {
				allCon3[i], allCon3[i+1] = allCon3[i+1], allCon3[i]
			}
		}
	}
	for _, id := range allCon3 {
		allCon = append(allCon, id.usID)
	}
	fmt.Println("check all cons", allCon3)
	// allCon = append(allCon, recID)
	// for i := 0; i < len(allCon)/2; i++ {
	// 	allCon[i], allCon[len(allCon)-(i+1)] = allCon[len(allCon)-(i+1)], allCon[i]
	// }
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
	fmt.Println("sorted cons", lastOne)
	return lastOne
}

func ChangeNotif(curUserID, senderID int) {
	var newArr []int
	notif := FindNotification(curUserID)
	newArr = append(newArr, notif...)
	for index, in := range newArr {
		if in == senderID {
			newArr = remove(newArr, index)
		}
	}
	slcNotif := make([]string, len(newArr))
	for i := 0; i < len(newArr); i++ {
		str := strconv.Itoa(newArr[i])
		slcNotif[i] = str
	}
	newNotificationString := strings.Join(slcNotif, ",")
	rows2, err := db.Prepare("UPDATE users SET notifications = ? WHERE userID = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()
	rows2.Exec(newNotificationString, curUserID)
	fmt.Println(len(newNotificationString), "notif string", newNotificationString)
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
