package forum

import (
	"encoding/json"
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

var pwHashDB []byte

func loginFailed(conn *websocket.Conn) {
	// login failed
	fmt.Println("Login Failed")
	var failedResponse WsLoginResponse
	failedResponse.Label = "login"
	failedResponse.Content = "Please check your credentials"
	failedResponse.Pass = false
	conn.WriteJSON(failedResponse)
	return
	// return false
}

func LoginWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Login Connected")
	var firstResponse WsLoginResponse
	firstResponse.Label = "greet"
	// firstResponse.Content = "Please login to the Forum"
	conn.WriteJSON(firstResponse)
	// insert conn into db with empty userID, fill in the userID when registered or logged in
	// stmt, err := db.Prepare(`INSERT INTO websockets (userID, websocketAdd) VALUES (?, ?);`)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// stmt.Exec("", conn)

	// loginSuccess := false
	// // keep running until login success
	// for !loginSuccess {
	// loginSuccess = listenToLoginWs(conn)
	// }
	// conn.Close()

	// if loginSuccess {
	// 	userIsOnline()
	// }

	listenToLoginWs(conn)
}

func listenToLoginWs(conn *websocket.Conn) {
	defer func() {
		fmt.Println("Login Ws Conn Closed")
	}()

	var loginPayload WsLoginPayload

	for {
		err := conn.ReadJSON(&loginPayload)
		if err == nil {
			// loginPayload.Conn = conn
			fmt.Printf("login payload received: %v\n", loginPayload)
			// testLogin() // just for testing, can be removed in production
			// loginSuccess := ProcessAndReplyLogin(conn, loginPayload)
			// return loginSuccess
			ProcessAndReplyLogin(conn, loginPayload)
		}
	}
}

func ProcessAndReplyLogin(conn *websocket.Conn, loginPayload WsLoginPayload) {
	fmt.Printf("login u: %s: , login pw: %s\n", loginPayload.NicknameEmail, loginPayload.Password)
	// // get user data from db
	var logge bool
	var logUser User
	var not string
	// auth user
	fmt.Printf("%s trying to Login\n", loginPayload.NicknameEmail)
	rows, err := db.Query(`SELECT *
							FROM users
							WHERE nickname = ?
							OR email = ?`, loginPayload.NicknameEmail, loginPayload.NicknameEmail)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&logUser.UserId, &logUser.Nickname, &logUser.Age, &logUser.Gender, &logUser.FirstName, &logUser.LastName, &logUser.Email, &pwHashDB, &logge, &not) // bug in db not is before logge?
	}
	if logUser.UserId == 0 {
		loginFailed(conn)
		return
	}
	findCurUser(logUser.UserId)

	// // test hash
	// hash, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	// fmt.Printf("nicknameEmailDB: %s , hashDB: %s\n", nicknameEmailDB, hashDB)

	// // compare pw
	err = bcrypt.CompareHashAndPassword(pwHashDB, []byte(loginPayload.Password))
	// fmt.Printf("DB pw: %s, entered: %s\n", hashDB, loginPayload.password)
	// fmt.Printf("DB pw: %s, entered: %s\n", hashDB, hash)

	// Login failed
	if err != nil {
		loginFailed(conn)
		return
	} else {
		// Login successfully
		fmt.Printf("%s (name from DB) Login successfully\n", loginPayload.NicknameEmail)
		// update login status in users
		currentUser, err := json.Marshal(logUser)
		if err != nil {
			log.Fatal(err)
		}
		var successResponse WsLoginResponse
		successResponse.Label = "login"
		// no need the form is closed after success
		// successResponse.Content = fmt.Sprintf("%s Login successfully", nicknameDB)
		successResponse.Content = string(currentUser)
		successResponse.Pass = true
		successResponse.Cookie = genCookie(conn, logUser.UserId)
		conn.WriteJSON(successResponse)
		return
	}
	// return true
}

// func testLogin() {
// 	stmt, err := db.Prepare("INSERT INTO users (userID, nickname, age, gender, firstname, lastname, email, password, loggedIn) VALUES (?,?,?,?,?,?,?,?,?);")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	testpw := "supersecret"
// 	testpwHash, err := bcrypt.GenerateFromPassword([]byte(testpw), 10)
// 	stmt.Exec(7, "doubleOh7", 42, 1, "James", "Bond", "secretagent@mi5.com", testpwHash, false)
// }
