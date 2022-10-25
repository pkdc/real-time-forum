package forum

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	uuid "github.com/satori/go.uuid"
// 	"golang.org/x/crypto/bcrypt"
// )

// type loginData struct {
// 	NicknameEmail string `json:"name"`
// 	Password      string `json:"pw"`
// }

// type loginResponse struct {
// 	Pass bool `json:"pass"`
// }

// func loggedIn(r *http.Request) bool {
// 	// c, err := r.Cookie("session")
// 	// if err != nil {
// 	// 	// check if the uuid exists in the session table
// 	// 	return false
// 	// }
// 	// var unameTryingToLogin string
// 	// var sid string
// 	// rows, err := db.Query("SELECT username, sessionID FROM sessions WHERE sessionID = ?;", c.Value)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// defer rows.Close()
// 	// for rows.Next() {
// 	// 	rows.Scan(&unameTryingToLogin, &sid)
// 	// }
// 	// if unameTryingToLogin == "" && sid == "" {
// 	return false
// 	// }
// 	// return true
// }

// func processLogin(w http.ResponseWriter, r *http.Request) {
// 	// get login data from request

// 	fmt.Printf("login u: %s: , login pw: %s\n", loginPayload.NicknameEmail, loginPayload.Password)

// 	// // get user data from db
// 	var nicknameEmailDB string
// 	var hashDB []byte

// 	fmt.Printf("%s trying to Login\n", loginPayload.NicknameEmail)
// 	rows, err := db.Query(`SELECT nickname, email, password
// 							FROM users
// 							WHERE nickname = ?
// 							OR email = ?`, loginPayload.NicknameEmail, loginPayload.NicknameEmail)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		rows.Scan(&nicknameEmailDB, &hashDB)
// 	}

// 	// // test hash
// 	// hash, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
// 	// fmt.Printf("nicknameEmailDB: %s , hashDB: %s\n", nicknameEmailDB, hashDB)

// 	// // compare pw
// 	err = bcrypt.CompareHashAndPassword(hashDB, []byte(loginPayload.Password))
// 	// fmt.Printf("DB pw: %s, entered: %s\n", hashDB, loginPayload.password)
// 	// fmt.Printf("DB pw: %s, entered: %s\n", hashDB, hash)
// 	if err != nil {
// 		// login failed
// 		var response loginResponse
// 		response.Pass = false
// 		js, err := json.Marshal(response)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(js)
// 		// http.Error(w, "Username or Password not matched, please try again", http.StatusForbidden)
// 		// tpl, err := template.ParseFiles("./templates/failedLogin.gohtml")
// 		// if err != nil {
// 		// 	http.Error(w, "Parsing Error", http.StatusInternalServerError)
// 		// }
// 		// tpl.Execute(w, nil)
// 		// return
// 	}
// 	fmt.Printf("%s (name from DB) Login successfully\n", loginPayload.NicknameEmail)

// 	// // allow each user to have only one opened session
// 	// var loggedInUname string
// 	// rows, err = db.Query("SELECT username FROM sessions WHERE username = ?;", nicknameEmailDB)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// defer rows.Close()
// 	// for rows.Next() {
// 	// 	rows.Scan(&loggedInUname)
// 	// }
// 	// // if the uname can be found in session table, remove that row (should only have 1 row)
// 	// if loggedInUname != "" {
// 	// 	stmt, err := db.Prepare("DELETE FROM sessions WHERE username = ?;")
// 	// 	if err != nil {
// 	// 		log.Fatal(err)
// 	// 	}
// 	// 	defer stmt.Close()
// 	// 	stmt.Exec(loggedInUname)
// 	// }

// 	// // assign a cookie
// 	sid := uuid.NewV4()
// 	fmt.Printf("login sid: %s\n", sid)
// 	http.SetCookie(w, &http.Cookie{
// 		Name:   "session",
// 		Value:  sid.String(),
// 		MaxAge: 900, // 15mins
// 	})

// 	// // forumUser.Username = nicknameEmailDB
// 	// // forumUser.LoggedIn = true
// 	// // fmt.Printf("%s forum User Login\n", forumUser.Username)

// 	// // update the user's login status
// 	stmt, err := db.Prepare("UPDATE users SET loggedIn = ? WHERE userID = ?;")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
// 	stmt.Exec(true, nicknameEmailDB)

// 	// // insert a record into session table
// 	stmt, err = db.Prepare("INSERT INTO sessions (sessionID, userID) VALUES (?, ?);")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
// 	stmt.Exec(sid.String(), nicknameEmailDB)

// 	//test
// 	// var whichUser string
// 	// var logInOrNot bool
// 	// rows, err = db.Query("SELECT username, loggedIn FROM users WHERE username = ?;", nicknameEmailDB)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// defer rows.Close()
// 	// for rows.Next() {
// 	// 	rows.Scan(&whichUser, &logInOrNot)
// 	// }
// 	// fmt.Printf("login user: %s, login status: %v\n", whichUser, logInOrNot)
// }

// // func processLogout(w http.ResponseWriter, r *http.Request) {
// // 	c, err := r.Cookie("session")
// // 	var logoutUname string

// // 	if err == nil {
// // 		// get the username of the logout user
// // 		rows, err := db.Query("SELECT username FROM sessions WHERE sessionID = ?;", c.Value)
// // 		if err != nil {
// // 			log.Fatal(err)
// // 		}
// // 		defer rows.Close()
// // 		for rows.Next() {
// // 			rows.Scan(&logoutUname)
// // 		}
// // 		fmt.Printf("Found user %s wants to logout", logoutUname)

// // 		// delete sessionID from sessions db table
// // 		stmt, err := db.Prepare("DELETE FROM sessions WHERE sessionID=?")
// // 		if err != nil {
// // 			log.Fatal(err)
// // 		}
// // 		defer stmt.Close()
// // 		stmt.Exec(c.Value)
// // 		fmt.Printf("cookie sid removed (have value): %s\n", c.Value)
// // 	}

// // 	// test
// // 	var sessionID string
// // 	rows, err := db.Query("SELECT * FROM sessions")
// // 	for rows.Next() {
// // 		rows.Scan(&sessionID)
// // 	}
// // 	fmt.Printf("cookie sid removed (should be empty): %s\n", sessionID) // empty is correct

// // 	// delete browser's cookie
// // 	_, err = r.Cookie("session")
// // 	if err == nil {
// // 		http.SetCookie(w, &http.Cookie{
// // 			Name:   "session",
// // 			Value:  "",
// // 			MaxAge: -1,
// // 		})
// // 	}
// // 	fmt.Printf("%s Logout\n", logoutUname)

// // 	stmt, err := db.Prepare("UPDATE users SET loggedIn = ? WHERE username = ?;")
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	defer stmt.Close()
// // 	stmt.Exec(false, logoutUname)
// // }

// // func obtainCurUserFormCookie(r *http.Request) user {
// // 	var curUser user
// // 	c, err := r.Cookie("session")
// // 	// if there is a session cookie
// // 	if err == nil {
// // 		fmt.Printf("There is a cookie, sid: %s\n", c.Value)
// // 		// get current username from session table
// // 		var curUname string
// // 		rows, err := db.Query("SELECT username FROM sessions WHERE sessionID = ?;", c.Value)
// // 		if err != nil {
// // 			log.Fatal(err)
// // 		}
// // 		defer rows.Close()
// // 		for rows.Next() {
// // 			rows.Scan(&curUname)
// // 			// fmt.Printf("Found uname %s in sessions\n", curUname)
// // 		}
// // 		fmt.Printf("Found uname %s in sessions\n", curUname)
// // 		rows, err = db.Query("SELECT username, image, email, access, loggedIN  FROM users WHERE username = ?;", curUname)
// // 		if err != nil {
// // 			log.Fatal(err)
// // 		}
// // 		defer rows.Close()
// // 		for rows.Next() {
// // 			rows.Scan(&curUser.Username, &curUser.Image, &curUser.Email, &curUser.Access, &curUser.LoggedIn)
// // 			fmt.Printf("Found user %s in users, with login status %v\n", curUser.Username, curUser.LoggedIn)
// // 		}
// // 	}

// // 	// test
// // 	var whichUser string
// // 	var logInOrNot bool
// // 	rows, err := db.Query("SELECT username, loggedIn FROM users WHERE username = ?;", curUser.Username)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	defer rows.Close()
// // 	for rows.Next() {
// // 		rows.Scan(&whichUser, &logInOrNot)
// // 	}
// // 	fmt.Printf("obtainCurUserFormCookie:: login user: %s, login status: %v\n", whichUser, logInOrNot)

// // 	return curUser
// // }

// // func GetOutboundIP() net.IP {
// // 	conn, err := net.Dial("udp", "8.8.8.8:80")
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	defer conn.Close()

// // 	localAddr := conn.LocalAddr().(*net.UDPAddr)

// // 	return localAddr.IP
// // }
