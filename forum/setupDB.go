package forum

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func createUsersTable() {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS users (
		userID INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname VARCHAR(30),
		age INTEGER,
		gender VARCHAR(10),
		firstname VARCHAR(30),
		lastname VARCHAR(30),
		email VARCHAR(50),
		password VARCHAR(100),
		loggedIn BOOLEAN,
		profilepicture VARCHAR(100),
		notifications VARCHAR(100));`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func createSessionsTable() {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS sessions 
	(sessionID VARCHAR(30) PRIMARY KEY,
		userID INTEGER,
		FOREIGN KEY(userID) REFERENCES users(userID));`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func createPostsTable() {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS posts 
	(postID INTEGER PRIMARY KEY AUTOINCREMENT,
		userID INTEGER,
		title VARCHAR(50),
		content VARCHAR(1000),
		category VARCHAR(50),
		postTime DATETIME,
		FOREIGN KEY(userID) REFERENCES users(userID));`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func createCommentsTable() {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments 
	(commentID INTEGER PRIMARY KEY AUTOINCREMENT, 
		userID INTEGER, 
		postID INTEGER, 
		content VARCHAR(2000), 
		comTime DATETIME,
		FOREIGN KEY(userID) REFERENCES users(userID),
		FOREIGN KEY(postID) REFERENCES posts(postID));`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func createMessageTable() {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS messages
	(messageID INTEGER PRIMARY KEY AUTOINCREMENT,
		senderID INTEGER,
		receiverID INTEGER,
		messageTIme VARCHAR(2000),
		content VARCHAR(2000),
		seen BOOLEAN,
		FOREIGN KEY(senderID) REFERENCES users(userID),
		FOREIGN KEY(receiverID) REFERENCES users(userID));`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

// func createWebsocketsTable() {
// 	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS websockets
// 	(websocketID INTEGER PRIMARY KEY AUTOINCREMENT,
// 		userID INTEGER,
// 		websocketAdd VARCHAR(2000),
// 		usage VARCHAR(10),
// 		FOREIGN KEY(userID) REFERENCES users(userID));`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
// 	stmt.Exec()
// }

func InitDB() {
	db, _ = sql.Open("sqlite3", "./forum.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	createSessionsTable()
	createUsersTable()
	createPostsTable()
	createCommentsTable()
	createMessageTable()

	// InsertMessage(1,4, "hello")
	// InsertMessage(1,2, "hello")
	// InsertMessage(2,1, "hello")
	// InsertMessage(1,2, "how are you")
	// InsertMessage(2,1, "thanks")
	// InsertMessage(2,1, "i am fine and you")

	// createWebsocketsTable()
}

func InsertMessage(usID, recID int, cont string) {
	rows, err := db.Prepare("INSERT INTO messages (senderID, receiverID, messageTime,content,seen) VALUES (?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	rows.Exec(usID, recID, time.Now(), cont, false)
}
