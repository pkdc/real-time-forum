package forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type WsRegisterResponse struct {
	Label   string        `json:"label"`
	Content string        `json:"content"`
	Pass    bool          `json:"pass"`
	Cookie  SessionCookie `json:"cookie"`
}

type WsRegisterPayload struct {
	Label     string `json:"label"`
	FirstName string `json:"name"`
	LastName  string `json:"lastname"`
	NickName  string `json:"nickname"`
	Age       string `json:"age"`
	Email     string `json:"email"`
	Password  string `json:"pw"`
	Gender    string `json:"gender_option"`
}

var userID int

func RegWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Reg Connected")
	var firstResponse WsLoginResponse
	firstResponse.Label = "greet"
	firstResponse.Content = "Please register to the Forum"
	conn.WriteJSON(firstResponse)

	regSuccess := false
	for !regSuccess {
		regSuccess = listenToRegWs(conn)
	}
}

func listenToRegWs(conn *websocket.Conn) bool {
	defer func() {
		fmt.Println("Reg Ws Conn Closed")
	}()

	var regPayload WsRegisterPayload

	for {
		err := conn.ReadJSON(&regPayload)
		if err == nil {
			fmt.Printf("payload received: %v\n", regPayload)
			regSuccess := ProcessAndReplyReg(conn, regPayload)
			return regSuccess
		}
	}
}

func ProcessAndReplyReg(conn *websocket.Conn, regPayload WsRegisterPayload) bool {
	dob, err := time.Parse("2006-01-02", regPayload.Age)
	if err != nil {
		log.Fatal(err)
	}
	year := time.Time.Year(dob)
	age := time.Now().Year() - year
	ageStr := strconv.Itoa(age)
	password := []byte(regPayload.Password)
	cryptPw, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		log.Fatal(err)
	}
	if regPayload.Label == "reg" {
		fmt.Printf("reg- FirstN: %s, LastN: %s, NickN : %s, age: %s, email %s, pw: %s, gender: %s\n",
			regPayload.FirstName, regPayload.LastName, regPayload.NickName,
			ageStr, regPayload.Email, cryptPw, regPayload.Gender)

		fmt.Printf("%s creating user\n", regPayload.NickName)
		rows, err := db.Prepare("INSERT INTO users(nickname,age,gender,firstname,lastname,email,password, loggedIn) VALUES(?,?,?,?,?,?,?,?);")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		rows.Exec(regPayload.NickName, ageStr, regPayload.Gender, regPayload.FirstName, regPayload.LastName, regPayload.Email, cryptPw, true)
		if regPayload.NickName != "" && ageStr != "" && regPayload.Gender != "" && regPayload.FirstName != "" && regPayload.LastName != "" && regPayload.Email != "" && cryptPw != nil {

			fmt.Println("Register successfully")

			var successResponse WsRegisterResponse
			successResponse.Label = "reg"
			successResponse.Content = fmt.Sprintf("%s Login successfully", regPayload.NickName)
			successResponse.Pass = true

			rows3, err := db.Query(`SELECT userID FROM users WHERE nickname = ?`, regPayload.NickName)
			if err != nil {
				log.Fatal(err)
			}
			defer rows3.Close()
			for rows3.Next() {
				rows3.Scan(&userID)
			}

			successResponse.Cookie = genCookie(conn, userID)
			conn.WriteJSON(successResponse)

		} else {
			var failedResponse WsRegisterResponse
			failedResponse.Label = "reg"
			failedResponse.Content = fmt.Sprintf("Check the credentials")
			failedResponse.Pass = false
			conn.WriteJSON(failedResponse)
			return false
		}
	}
	return true
}
