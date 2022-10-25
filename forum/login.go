package forum

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type SessionCookie struct {
	Uid    int    `json:"uid"`
	Sid    string `json:"sid"`
	MaxAge int    `json:"max_age"`
}

type WsLoginResponse struct {
	Label   string        `json:"label"`
	Content string        `json:"content"`
	Pass    bool          `json:"pass"`
	Cookie  SessionCookie `json:"cookie"`
}

type WsLoginPayload struct {
	Label         string `json:"label"`
	NicknameEmail string `json:"name"`
	Password      string `json:"pw"`
	// Conn          *websocket.Conn `json:"-"`
}

func LoginWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Connected")
	var firstResponse WsLoginResponse
	firstResponse.Label = "Greet"
	firstResponse.Content = "Please login to the Forum"
	conn.WriteJSON(firstResponse)
	// insert conn into db with empty userID, fill in the userID when registered or logged in
	// stmt, err := db.Prepare(`INSERT INTO websockets (userID, websocketAdd) VALUES (?, ?);`)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// stmt.Exec("", conn)
	loginSuccess := false

	// keep running until login success
	for !loginSuccess {
		loginSuccess = listenToLoginWs(conn)
	}
}

func listenToLoginWs(conn *websocket.Conn) bool {
	defer func() {
		fmt.Println("Ws Conn Closed")
	}()

	var loginPayload WsLoginPayload

	for {
		err := conn.ReadJSON(&loginPayload)
		if err == nil {
			// loginPayload.Conn = conn
			fmt.Printf("payload received: %v\n", loginPayload)
			testLogin() // just for testing, can be removed in production
			loginSuccess := ProcessAndReplyLogin(conn, loginPayload)
			return loginSuccess
		}
	}
}

func ProcessAndReplyLogin(conn *websocket.Conn, loginPayload WsLoginPayload) bool {

	fmt.Printf("login u: %s: , login pw: %s\n", loginPayload.NicknameEmail, loginPayload.Password)

	// // get user data from db
	var userIDDB int
	var nicknameDB string
	var emailDB string
	var hashDB []byte

	fmt.Printf("%s trying to Login\n", loginPayload.NicknameEmail)
	rows, err := db.Query(`SELECT userID, nickname, email, password 
							FROM users
							WHERE nickname = ?
							OR email = ?`, loginPayload.NicknameEmail, loginPayload.NicknameEmail)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&userIDDB, &nicknameDB, &emailDB, &hashDB)
	}

	// // test hash
	// hash, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	// fmt.Printf("nicknameEmailDB: %s , hashDB: %s\n", nicknameEmailDB, hashDB)

	// // compare pw
	err = bcrypt.CompareHashAndPassword(hashDB, []byte(loginPayload.Password))
	// fmt.Printf("DB pw: %s, entered: %s\n", hashDB, loginPayload.password)
	// fmt.Printf("DB pw: %s, entered: %s\n", hashDB, hash)

	// Login failed
	if err != nil {
		// login failed
		fmt.Println("Failed")
		var failedResponse WsLoginResponse
		failedResponse.Label = "login"
		failedResponse.Content = "record cannot be found"
		failedResponse.Pass = false
		conn.WriteJSON(failedResponse)
		return false
	}
	// Login successfully
	fmt.Printf("%s (name from DB) Login successfully\n", loginPayload.NicknameEmail)

	var successResponse WsLoginResponse
	successResponse.Label = "login"
	successResponse.Content = fmt.Sprintf("%s Login successfully", nicknameDB)
	successResponse.Pass = true
	successResponse.Cookie = genCookie(conn, userIDDB)
	conn.WriteJSON(successResponse)

	return true
}

func testLogin() {
	stmt, err := db.Prepare("INSERT INTO users (userID, nickname, age, gender, firstname, lastname, email, password, loggedIn) VALUES (?,?,?,?,?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	testpw := "supersecret"
	testpwHash, err := bcrypt.GenerateFromPassword([]byte(testpw), 10)
	stmt.Exec(007, "doubleOh7", 42, 1, "James", "Bond", "secretagent@mi5.com", testpwHash, false)
}
