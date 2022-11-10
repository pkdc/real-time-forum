package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WsPostResponse struct {
	Label   string `json:"label"`
	Content string `json:"content"`
	Pass    bool   `json:"pass"`
}

type Ind struct {
	Index int           `json:"index"`
	Post  WsPostPayload `json:"postinfo"`
}

type WsPostPayload struct {
	Label    string `json:"label"`
	UserID   string `json:"userID"`
	Title    string `json:"title"`
	Category string `json:"category_option"`
	Content  string `json:"Content"`
	PostTime string `json:"PostTime"`
}

func findAllPosts() string {
	var pos []WsPostPayload
	var everyPost []Ind
	rows, err := db.Query("SELECT title,content,category,postTime FROM posts GROUP BY postID;")
	if err != nil {
		log.Fatal(err)
	}
	// ----------------------- DONT FORGET --------------------
	// after session done, i will add userID there.
	defer rows.Close()
	for rows.Next() {
		var po WsPostPayload
		var postTime time.Time
		rows.Scan(&(po.Title), &(po.Content), &(po.Category), &postTime)
		po.PostTime = postTime.Format("Mon 02-01-2006 15:04:05")
		fmt.Println(postTime)
		pos = append(pos, po)
		fmt.Println("THIS IS POST", po)
	}
	for i := 0; i < len(pos); i++ {
		var singlePost Ind
		singlePost.Index = i
		singlePost.Post = pos[i]
		everyPost = append(everyPost, singlePost)
	}
	j, err := json.Marshal(everyPost)
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}

func PostWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Post Connected")
	var firstResponse WsPostResponse
	firstResponse.Label = "Greet"
	allPosts := findAllPosts()
	firstResponse.Content = allPosts
	conn.WriteJSON(firstResponse)
	listenToPostWs(conn)
}

func listenToPostWs(conn *websocket.Conn) {
	defer func() {
		fmt.Println("Post Ws Conn Closed")
	}()

	var postPayload WsPostPayload

	for {
		err := conn.ReadJSON(&postPayload)
		if err == nil {
			fmt.Printf("payload received: %v\n", postPayload)
			ProcessAndReplyPost(conn, postPayload)
		}
	}
}

func ProcessAndReplyPost(conn *websocket.Conn, postPayload WsPostPayload) {
	if postPayload.Label == "post" {
		fmt.Println("LABEL WORK--------------------------------")
		fmt.Printf("post - title:%s, cat:%s, Content:%s", postPayload.Title, postPayload.Category, postPayload.Content)

		rows, err := db.Prepare("INSERT INTO posts(title,content,category,postTime) VALUES(?,?,?,?);")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		rows.Exec(postPayload.Title, postPayload.Content, postPayload.Category, time.Now())
		fmt.Println("Posted successfully")
		var successResponse WsPostResponse
		successResponse.Label = "post"
		successResponse.Content = findAllPosts()
		successResponse.Pass = true
		conn.WriteJSON(successResponse)

	}
}
