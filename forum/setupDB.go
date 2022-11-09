package forum

import (
	"database/sql"
	"log"

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
		loggedIn BOOLEAN);`)
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
		commentTime DATETIME, 
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
		messageTIme DATETIME,
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

func createWebsocketsTable() {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS websockets
	(websocketID INTEGER PRIMARY KEY AUTOINCREMENT,
		userID INTEGER,
		websocketAdd VARCHAR(2000),
		usage VARCHAR(10),
		FOREIGN KEY(userID) REFERENCES users(userID));`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

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
	createWebsocketsTable()
}
