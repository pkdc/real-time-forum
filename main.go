package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	InitDB()
	go ProcessAndReplyUserList()
	hub := NewHub()
	hub.Run()
	// ClearUsers()
	// ClearPosts()
	// ClearComments()
	exec.Command("xdg-open", "http://localhost:8080/").Start()
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/regWs/", RegisterHandler)
	http.HandleFunc("/postWs/", PostWsEndpoint)
	http.HandleFunc("/loginWs/", LoginHandler)
	http.HandleFunc("/userListWs/", UserListHandler)
	http.HandleFunc("/chatWs/", ChatHandler)

	// http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/logout/", LogoutHandler)
	// http.HandleFunc("/postpage", PostPageHandler)
	// http.HandleFunc("/notifications", NotiPageHandler)
	// http.HandleFunc("/activity", ActivityPageHandler)
	// // http.HandleFunc("/delete", DeleteHandler)
	fmt.Println("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
