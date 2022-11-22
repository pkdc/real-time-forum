package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("./assets/index.html")
	if err != nil {
		http.Error(w, "Parsing Error", http.StatusInternalServerError)
		return
	}
	err = tpl.ExecuteTemplate(w, "index.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("logged in", loggedInCheck(r))
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	if r.Method == "GET" && !loggedInCheck(r) {
		LoginWsEndpoint(w, r)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	if r.Method == http.MethodGet && !loggedInCheck(r) {
		RegWsEndpoint(w, r)
	}
}

func UserListHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("logged in", loggedInCheck(r))
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	if r.Method == "GET" && !loggedInCheck(r) {
		userListWsEndpoint(w, r)
	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("logged in", loggedInCheck(r))
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	if r.Method == "GET" && !loggedInCheck(r) {
		chatWsEndpoint(w, r)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if loggedInCheck(r) {
		processLogout(w, r)
	}
}
