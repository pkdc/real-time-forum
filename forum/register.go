package forum

import (
	"encoding/json"
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
	ProfilePicture string `json:"pp_option"`
}

type User struct {
	UserId    int    `json:"userID"`
	Nickname  string `json:"nickname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	ProfilePicture string `json:"pp"`
	LoggedIn  bool
}

var (
	userID  int
	curUser User
)

func findCurUser(userid int) {
	// var id int
	// var pw string
	// var curNickname string
	// var curAge int
	// var curGender string
	// var curFirstName string
	// var curLastName string
	// var curEmail string
	// var curLoggedIn bool
	rows, err := db.Query(`SELECT nickname, age, gender, firstname, lastname, email,loggedIn,profilepicture FROM users WHERE userID = ?`, userid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		// rows3.Scan(&id, &curNickname, &curAge, &curGender, &curFirstName, &curLastName, &curEmail, &pw, &curLoggedIn)
		rows.Scan(&curUser.Nickname, &curUser.Age, &curUser.Gender, &curUser.FirstName, &curUser.LastName, &curUser.Email, &curUser.LoggedIn, &curUser.ProfilePicture)
	}
	// curUser.Nickname = curNickname
	// curUser.Age = curAge
	// curUser.Gender = curGender
	// curUser.FirstName = curFirstName
	// curUser.LastName = curLastName
	// curUser.Email = curEmail
	// curUser.LoggedIn = curLoggedIn

	// fmt.Printf("nickname: %v \n", curUser.Nickname)
}

func RegWsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Reg Connected")
	var firstResponse WsLoginResponse
	firstResponse.Label = "greet"
	// firstResponse.Content = "Please register to the Forum"
	conn.WriteJSON(firstResponse)

	// regSuccess := false
	// for !regSuccess {
	// 	regSuccess = listenToRegWs(conn)
	// }
	listenToRegWs(conn)
}

func listenToRegWs(conn *websocket.Conn) {
	defer func() {
		fmt.Println("Reg Ws Conn Closed")
	}()

	var regPayload WsRegisterPayload

	for {
		err := conn.ReadJSON(&regPayload)
		if err == nil {
			fmt.Printf("payload received: %v\n", regPayload)
			// regSuccess := ProcessAndReplyReg(conn, regPayload)
			// return regSuccess
			ProcessAndReplyReg(conn, regPayload)
		}
	}
}

func ProcessAndReplyReg(conn *websocket.Conn, regPayload WsRegisterPayload) {
	var emailCheck string
	var nicknameCheck string
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
		// checking duplicate
		rows2, err := db.Query(`SELECT nickname FROM users WHERE email = ?`, regPayload.Email)
		if err != nil {
			log.Fatal(err)
			// return false
		}
		defer rows2.Close()
		for rows2.Next() {
			rows2.Scan(&emailCheck)
			if len(emailCheck) != 0 {
				fmt.Println("already registered")
				var failedResponse WsRegisterResponse
				failedResponse.Label = "reg"
				failedResponse.Content = "ERROR - This email already taken"
				failedResponse.Pass = false
				conn.WriteJSON(failedResponse)
				return
				// return false
			}
		}
		rows5, err := db.Query(`SELECT email FROM users WHERE nickname = ?`, regPayload.NickName)
		if err != nil {
			log.Fatal(err)
			// return false
		}
		defer rows5.Close()
		for rows5.Next() {
			rows5.Scan(&nicknameCheck)
			if len(nicknameCheck) != 0 {
				fmt.Println("already registered")
				var failedResponse WsRegisterResponse
				failedResponse.Label = "reg"
				failedResponse.Content = "ERROR - This nickname already taken"
				failedResponse.Pass = false
				conn.WriteJSON(failedResponse)
				return
				// return false
			}
		}
		// insert newuser  into database
		fmt.Printf("%s creating user\n", regPayload.NickName)
		stmt, err := db.Prepare("INSERT INTO users(nickname,age,gender,firstname,lastname,email,password, loggedIn, profilepicture) VALUES(?,?,?,?,?,?,?,?,?);")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(regPayload.NickName, ageStr, regPayload.Gender, regPayload.FirstName, regPayload.LastName, regPayload.Email, cryptPw, false,regPayload.ProfilePicture)

		if regPayload.NickName != "" && ageStr != "" && regPayload.Gender != "" && regPayload.FirstName != "" && regPayload.LastName != "" && regPayload.Email != "" && cryptPw != nil {

			fmt.Println("Register successfully")

			var successResponse WsRegisterResponse
			successResponse.Label = "reg"
			// no need the form is closed after success
			// successResponse.Content = fmt.Sprintf("%s Login successfully", regPayload.NickName)
			successResponse.Pass = true

			rows3, err := db.Query(`SELECT userID FROM users WHERE nickname = ?`, regPayload.NickName)
			if err != nil {
				log.Fatal(err)
			}
			defer rows3.Close()
			for rows3.Next() {
				rows3.Scan(&userID)
			}
			findCurUser(userID)
			curUser.UserId = userID
			fmt.Println("curuser:", curUser)
			userJson, err := json.Marshal(curUser)
			if err != nil {
				log.Fatal(err)
			}
			successResponse.Content = string(userJson)
			// successResponse.Cookie = genCookie(conn, userID)
			conn.WriteJSON(successResponse)
		} else {
			var failedResponse WsRegisterResponse
			failedResponse.Label = "reg"
			failedResponse.Content = "Please check your credentials"
			failedResponse.Pass = false
			conn.WriteJSON(failedResponse)
			// return false
		}
	}
	// return true
}
