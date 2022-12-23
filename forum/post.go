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

type WsPostResponse struct {
	Label   string `json:"label"`
	Content string `json:"content"`
	Pass    bool   `json:"pass"`
}

type Ind struct {
	Index int           `json:"index"`
	Post  WsPostPayload `json:"postinfo"`
}

type IndCom struct {
	Index       int          `json:"indexCom"`
	CommentInfo WsComPayload `json:"comInfo"`
}

type WsPostPayload struct {
	Label         string `json:"label"`
	UserID        string `json:"userID"`
	Title         string `json:"title"`
	Category      string `json:"category_option"`
	Content       string `json:"Content"`
	PostTime      string `json:"PostTime"`
	PostID        string `json:"postID"`
	CommentOfPost string `json:"commentOfPost"`
	Conn          websocket.Conn
}

type WsComPayload struct {
	Label   string `json:"label"`
	UserID  string `json:"userID"`
	Comment string `json:"comment"`
	ComTime string `json:"comTime"`
}

var postComWsArr []*websocket.Conn
var lastCon *websocket.Conn

func findAllPosts() []Ind {
	var pos []WsPostPayload
	var everyPost []Ind
	var id int
	rows, err := db.Query("SELECT postID,title,content,category,postTime FROM posts GROUP BY postID;")
	if err != nil {
		log.Fatal(err)
	}
	// ----------------------- DONT FORGET --------------------
	// after session done, i will add userID there.
	defer rows.Close()
	for rows.Next() {
		var po WsPostPayload
		var postTime time.Time
		rows.Scan(&id, &(po.Title), &(po.Content), &(po.Category), &postTime)
		po.PostTime = postTime.Format("Mon 02-01-2006 15:04:05")
		po.CommentOfPost = findAllComments(id)
		pos = append(pos, po)
	}
	for i := 0; i < len(pos); i++ {
		var singlePost Ind
		singlePost.Index = i
		strI := strconv.Itoa(i + 1)
		pos[i].PostID = strI
		singlePost.Post = pos[i]
		everyPost = append(everyPost, singlePost)
	}
	// j, err := json.Marshal(everyPost)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// return string(j)
	return everyPost
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
	firstJson, err := json.Marshal(allPosts)
	if err != nil {
		log.Fatal(err)
	}
	firstResponse.Content = string(firstJson)
	conn.WriteJSON(firstResponse)
	lastCon= conn
	postComWsArr = append(postComWsArr, conn)
	ListenToPostWs()
	
}

func ListenToPostWs() {
	defer func() {
		fmt.Println("Post Ws Conn Closed")
	}()
	var postPayload WsPostPayload
	for {
		err := lastCon.ReadJSON(&postPayload)
		if err == nil {

			if postPayload.Label == "post" {
				fmt.Printf("payload received: %v\n", postPayload)
				ProcessAndReplyPost(lastCon, postPayload)
			} else if postPayload.Label == "Createcomment" {		
				ProcessAndReplyPost(lastCon, postPayload)
			} else if postPayload.Label == "showComment" {
				ProcessAndReplyPost(lastCon, postPayload)
			}
		}
	}
}

func ProcessAndReplyPost(conn *websocket.Conn, postPayload WsPostPayload) {
	if postPayload.Label == "post" {
		fmt.Printf("post - title:%s, cat:%s, Content:%s", postPayload.Title, postPayload.Category, postPayload.Content)

		rows, err := db.Prepare("INSERT INTO posts(title,content,category,postTime) VALUES(?,?,?,?);")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		rows.Exec(postPayload.Title, postPayload.Content, postPayload.Category, time.Now())
		fmt.Println("Post saved successfully")
		var successResponse WsPostResponse
		successResponse.Label = "post"
		allPosts := findAllPosts()
		postJson, err := json.Marshal(allPosts)
		if err != nil {
			log.Fatal(err)
		}
		successResponse.Content = string(postJson)
		successResponse.Pass = true
		// conn.WriteJSON(successResponse)
		broadcastPost(successResponse)

	} else if postPayload.Label == "Createcomment" {
		rows, err := db.Prepare("INSERT INTO comments (content, postID, comTime) VALUES (?,?,?);")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var comUnmars WsComPayload
		json.Unmarshal([]byte(postPayload.CommentOfPost), &comUnmars)
		comTime := time.Now()
		rows.Exec(comUnmars.Comment, postPayload.PostID, comTime)
		fmt.Println("comment saved successfully")
		var successResponse WsPostResponse
		successResponse.Label = "Createcomment"
		allPosts := findAllPosts()
		postJson, err := json.Marshal(allPosts)
		if err != nil {
			log.Fatal(err)
		}
		successResponse.Content = string(postJson)
		successResponse.Pass = true
		// conn.WriteJSON(successResponse)
		broadcastPost(successResponse)
	} else if postPayload.Label == "showComment" {
		var successResponse WsPostResponse
		successResponse.Label = "showComment"
		allPosts := findAllPosts()
		postJson, err := json.Marshal(allPosts)
		if err != nil {
			log.Fatal(err)
		}
		successResponse.Content = string(postJson)
		successResponse.Pass = true
		// conn.WriteJSON(successResponse)
		broadcastPost(successResponse)
	}
}

func findAllComments(postID int) string {
	// postID, err := strconv.Atoi(postIDstr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var com []WsComPayload
	var everyCom []IndCom
	var timeofCom time.Time
	rows, err := db.Query("SELECT content, comTime FROM comments WHERE postID = ?;", postID)
	if err != nil {
		log.Fatal(err)
	}
	// ----------------------- DONT FORGET --------------------
	// after session done, i will add userID there.
	defer rows.Close()
	for rows.Next() {
		var co WsComPayload
		rows.Scan(&(co.Comment), timeofCom)
		co.ComTime = timeofCom.Format("Mon 02-01-2006 15:04:05")
		com = append(com, co)
	}
	for i := 0; i < len(com); i++ {
		var singleCom IndCom
		singleCom.Index = i
		singleCom.CommentInfo = com[i]
		everyCom = append(everyCom, singleCom)
	}
	j, err := json.Marshal(everyCom)
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}

func broadcastPost(response WsPostResponse) {
	for _, postConns := range postComWsArr {
		postConns.WriteJSON(response)
	}
}
