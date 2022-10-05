package forum

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func createUsersTable() {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS users (
		userID INTEGER PRIMARY KEY,
		nickname VARCHAR(30),
		age INTEGER,
		gender BOOLEAN,
		firstname VARCHAR(30),
		lastname VARCHAR(30),
		image VARCHAR(2083),
		email VARCHAR(50),
		password VARCHAR(100),
		loggedIn BOOLEAN,
		notifyview VARCHAR(100),
		notifymsg VARCHAR(100),
		likedComments VARCHAR(100));`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func createSessionsTable() {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS sessions 
	(sessionID VARCHAR(50) PRIMARY KEY,
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
		author VARCHAR(30),image VARCHAR(2083),
		title VARCHAR(50),
		content VARCHAR(1000),
		category VARCHAR(50),
		postTime DATETIME,
		ips VARCHAR(10),
		URL VARCHAR(10),
		deleted BOOLEAN ,
		FOREIGN KEY(author) REFERENCES users(userID));`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec()
}

func createCommentsTable() {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments 
	(commentID INTEGER PRIMARY KEY AUTOINCREMENT, 
		author VARCHAR(30), 
		postID INTEGER, 
		content VARCHAR(2000), 
		commentTime DATETIME, 
		URL VARCHAR(10),
		FOREIGN KEY(author) REFERENCES users(author),
		FOREIGN KEY(postID) REFERENCES posts(postID));`)
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
}
