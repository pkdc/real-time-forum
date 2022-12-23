package forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WsChatResponse struct {
	Label     string `json:"label"`
	Content   string `json:"content"`
	UserID    int    `json:"userID"` // sender
	Sender    string `json:"sender"`
	ContactID int    `json:"contactID"` // receiver
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
	chatWsMap       = make(map[int]*websocket.Conn)
)

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

		if err == nil && chatPayload.Label == "chat" {
			processMsg(chatPayload)
			if userListWsMap[chatPayload.ReceiverId] != nil {
				chatPayloadChan <- chatPayload
			}
		} else if err == nil && chatPayload.Label == "updateChat" {
			// saving websocket to map
			chatWsMap[chatPayload.SenderId] = conn
		} else if err == nil && chatPayload.Label == "typing" {
			fmt.Printf("typing: chatPayload sender id: %d\n", chatPayload.SenderId)
			fmt.Printf("typing: chatPayload ReceiverId: %d\n", chatPayload.ReceiverId)
			// receiver is online
			if userListWsMap[chatPayload.ReceiverId] != nil {
				chatPayloadChan <- chatPayload
			}
		}
	}
}

func ProcessAndReplyChat() {
	for {
		receivedChatPayload := <-chatPayloadChan

		var responseChatPayload WsChatResponse

		if receivedChatPayload.Label == "chat" {
			responseChatPayload.Label = "msgIncoming"
			responseChatPayload.UserID = receivedChatPayload.SenderId
			responseChatPayload.ContactID = receivedChatPayload.ReceiverId
			responseChatPayload.Content = receivedChatPayload.Content
			receiverConn := chatWsMap[receivedChatPayload.ReceiverId]
			err := receiverConn.WriteJSON(responseChatPayload)
			updateUList()
			if err != nil {
				fmt.Println("failed to send message")
			}
		} else if receivedChatPayload.Label == "typing" {
			findCurUser(receivedChatPayload.SenderId)
			responseChatPayload.Label = "sender-typing"
			responseChatPayload.UserID = receivedChatPayload.SenderId
			responseChatPayload.ContactID = receivedChatPayload.ReceiverId
			responseChatPayload.Sender = curUser.Nickname
			fmt.Printf("typing: responseChatPayload sender id: %d\n", responseChatPayload.UserID)
			fmt.Printf("typing: responseChatPayload sender name: %s\n", responseChatPayload.Sender)
			fmt.Printf("typing: responseChatPayload ReceiverId: %d\n", responseChatPayload.ContactID)
			receiverConn := chatWsMap[receivedChatPayload.ReceiverId]
			err := receiverConn.WriteJSON(responseChatPayload)
			if err != nil {
				fmt.Println("failed to send message")
			}
		}

	}
}

func processMsg(msg WsChatPayload) {
	newNotif := true
	fmt.Println("msg:", msg)
	rows, err := db.Prepare("INSERT INTO messages(senderID,receiverID,messageTime,content,seen) VALUES(?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	rows.Exec(msg.SenderId, msg.ReceiverId, time.Now(), msg.Content, false)
	fmt.Println("msg saved successfully")
	notif := FindNotification(msg.ReceiverId)
	fmt.Println("OLD NOTIFICATION ARRAY", notif)
	if notif != nil {
		for _, not := range notif {
			if not == msg.SenderId {
				newNotif = false
				break
			}
		}
	}
	fmt.Println("notification bool:", newNotif)
	if newNotif {
		var newArr []int
		newArr = append(newArr, notif...)
		newArr = append(newArr, msg.SenderId)
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
		rows2.Exec(newNotificationString, msg.ReceiverId)
		fmt.Println(len(newNotificationString), "notif string", newNotificationString)
	}
}

func FindNotification(userID int) []int {
	notifi := ""
	rows, err := db.Query("SELECT notifications FROM users WHERE userID =?", userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&notifi)
	}
	fmt.Println("notifi:", notifi)
	if notifi == "" {
		fmt.Println("CANT FOUND NOTIFICATION")
		return nil
	}
	arr := strings.Split(notifi, ",")
	in := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		integ, err := strconv.Atoi(arr[i])
		if err != nil {
			log.Fatal(err)
		}
		in[i] = integ

	}
	fmt.Println("array of int", in)
	return in
}
