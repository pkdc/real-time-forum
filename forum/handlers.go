package forum

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

	// return
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(js)
	// }
	if r.Method == "GET" && !loggedInCheck(r) {
		LoginWsEndpoint(w, r)
	}
	// if r.Method == http.MethodPost {
		// 	fmt.Printf("----login-POST-----\n")
		// 	// processLogin(w, r)
		// }
	}
	
	func RegisterHandler(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		if r.Method == http.MethodGet && !loggedInCheck(r) {
			RegWsEndpoint(w,r)
			}
		}
	
// // func HomeHandler(w http.ResponseWriter, r *http.Request) {
// // 	if r.URL.Path != "/" {
// // 		http.Error(w, "404 Page Not Found", 404)
// // 		return
// // 	}
// // 	allForumUnames := allForumUnames()
// // 	if r.Method != http.MethodGet && r.Method != http.MethodPost {
// // 		http.Error(w, "Bad request", http.StatusBadRequest)
// // 	}

// // 	u, _ := url.Parse(r.URL.String())

// // 	query := u.Query()

// // 	category := []string{"Blockchain", "AI", "Cybersecurity", "Mobile Development", "Videogames"}
// // 	query.Get("category-filter")
// // 	if len(query) != 0 {
// // 		var badrequest bool = false
// // 		for i := 0; i < len(category); i++ {
// // 			if query.Get("category-filter") == category[i] {
// // 				badrequest = true
// // 			}
// // 		}
// // 		for i := 0; i < len(allForumUnames); i++ {
// // 			if query.Get("author-filter") == allForumUnames[i] {
// // 				badrequest = true
// // 			}
// // 		}
// // 		if query.Get("liked-post") == "liked-post" {
// // 			badrequest = true
// // 		}
// // 		if !badrequest {
// // 			http.Error(w, "400 Bad Request", 400)
// // 			return
// // 		}
// // 	}
// // 	if r.Method != http.MethodGet && r.Method != http.MethodPost {
// // 		http.Error(w, "Bad request", http.StatusBadRequest)
// // 	}

// // 	changingPos = false
// // 	curUser := obtainCurUserFormCookie(r)
// // 	if curUser.Username != "" {
// // 		users := AllForumUsers()
// // 		for i := 0; i < len(users); i++ {
// // 			if users[i].Username == curUser.Username {
// // 				curUser.LikedPost = users[i].LikedPost
// // 				curUser.DislikedPost = users[i].DislikedPost
// // 				curUser.DislikedComments2 = users[i].DislikedComments2
// // 				curUser.LikedComments2 = users[i].LikedComments2
// // 				curUser.Notifymsg = users[i].Notifymsg
// // 				curUser.Notifyview = users[i].Notifyview
// // 				var something []string
// // 				curUser, something = UpdateNotif(curUser)
// // 				fmt.Println(something)
// // 				fmt.Println("NOOOOOOTIFMESAAGE", curUser.NotifMessageShow)
// // 			}
// // 		}
// // 		changingPos = true
// // 	}

// // 	// // test
// // 	// var whichUser string
// // 	// var logInOrNot bool
// // 	// rows, err := db.Query("SELECT username, loggedIn FROM users WHERE username = ?;", curUser.Username)
// // 	// if err != nil {
// // 	// 	log.Fatal(err)
// // 	// }
// // 	// defer rows.Close()
// // 	// for rows.Next() {
// // 	// 	rows.Scan(&whichUser, &logInOrNot)
// // 	// }

// // 	// fmt.Printf("HomeHandler:: login user: %s, login status: %v\n", whichUser, logInOrNot)
// // 	if r.Method == http.MethodGet {
// // 		w.Header().Set("Content-Type", "text/html; charset=utf-8")

// // 		tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/header2.gohtml", "./templates/footer.gohtml", "./templates/index.gohtml", "./templates/index2.gohtml")
// // 		// tpl, err := template.ParseFiles("./templates/index.gohtml")
// // 		if err != nil {
// // 			http.Error(w, "Parsing Error", http.StatusInternalServerError)
// // 			return
// // 		}

// // 		filCat := r.FormValue("category-filter")
// // 		filAuthor := r.FormValue("author-filter")
// // 		filLiked := r.FormValue("liked-post")
// // 		filCatFromButton := r.FormValue("categoryOfPost")
// // 		var pos []post
// // 		if filCat != "" {
// // 			pos = filCatDisplayPostsAndComments(filCat)
// // 		} else if filAuthor != "" {
// // 			pos = filAuthorDisplayPostsAndComments(filAuthor)
// // 		} else if filLiked != "" {
// // 			pos = filLikedDisplayPostsAndComments(curUser)
// // 		} else if filCatFromButton != "" {
// // 			pos = filCatDisplayPostsAndComments(filCatFromButton)
// // 		} else {
// // 			pos = displayPostsAndComments()
// // 		}
// // 		AllLikes, AllDislikes := SumOfAllLikes(AllForumUsers())
// // 		pos = DistLikesToPosts(pos, AllLikes, AllDislikes)
// // 		for i := 0; i < len(pos); i++ {
// // 		}
// // 		if changingPos {

// // 			userLikes := CountLikesByUser(curUser, "l")
// // 			userDislikes := CountLikesByUser(curUser, "d")

// // 			for i := 0; i < len(pos); i++ {
// // 				for k := 0; k < len(userLikes); k++ {
// // 					if pos[i].PostID == userLikes[k] {
// // 						pos[i].LikedByCur = true
// // 					}
// // 				}
// // 			}
// // 			for i := 0; i < len(pos); i++ {
// // 				for k := 0; k < len(userDislikes); k++ {
// // 					if pos[i].PostID == userDislikes[k] {
// // 						pos[i].DislikedByCur = true
// // 					}
// // 				}
// // 			}
// // 		}

// // 		data := mainPageData{
// // 			Posts:       pos,
// // 			Userinfo:    curUser,
// // 			ForumUnames: allForumUnames,
// // 		}

// // 		// fmt.Println("---------", forumUser)
// // 		if changingPos {
// // 			err = tpl.ExecuteTemplate(w, "index2.gohtml", data)
// // 			if err != nil {
// // 				http.Error(w, "Executing Error", http.StatusInternalServerError)
// // 				return
// // 			}
// // 		} else {
// // 			err = tpl.ExecuteTemplate(w, "index.gohtml", data)
// // 			if err != nil {
// // 				http.Error(w, "Executing Error", http.StatusInternalServerError)
// // 				return
// // 			}
// // 		}

// // 	}
// // 	if r.Method == http.MethodPost {
// // 		processPost(r, curUser)
// // 		processComment(r, curUser)
// // 		http.Redirect(w, r, "/", http.StatusSeeOther)
// // 	}
// // }

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("logged in", loggedInCheck(r))
// 	if r.Method != http.MethodGet && r.Method != http.MethodPost {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 	}
// 	if loggedInCheck(r) {
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}
// 	if r.Method == "GET" {
// 		tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/footer.gohtml", "./templates/login.gohtml")
// 		if err != nil {
// 			http.Error(w, "Parsing Error", http.StatusInternalServerError)
// 			return
// 		}
// 		err = tpl.ExecuteTemplate(w, "login.gohtml", nil)
// 		if err != nil {
// 			http.Error(w, "Executing Error", http.StatusInternalServerError)
// 			return
// 		}
// 	}
// 	if r.Method == http.MethodPost {
// 		processLogin(w, r)
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 	}
// }
	// if r.Method == http.MethodPost {
	// 	regNewUser(w, r)
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// }
// }

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if loggedInCheck(r) {
		processLogout(w, r)
	}
}

// func PostPageHandler(w http.ResponseWriter, r *http.Request) {
// 	var strID string
// 	changingCom = false
// 	curUser := obtainCurUserFormCookie(r)
// 	if curUser.Username != "" {
// 		users := AllForumUsers()
// 		for i := 0; i < len(users); i++ {
// 			if users[i].Username == curUser.Username {
// 				curUser.LikedPost = users[i].LikedPost
// 				curUser.DislikedPost = users[i].DislikedPost
// 				curUser.DislikedComments2 = users[i].DislikedComments2
// 				curUser.LikedComments2 = users[i].LikedComments2

// 			}
// 		}
// 		changingCom = true

// 	}
// 	if r.Method == "GET" {
// 		tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/footer.gohtml", "./templates/header2.gohtml", "./templates/post.gohtml", "./templates/post2.gohtml")
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Parsing Error", http.StatusInternalServerError)
// 			return
// 		}
// 		strID = r.FormValue("postdetails")
// 		PostIdFromHTML, err := strconv.Atoi(strID)
// 		if err != nil {
// 			os.Exit(0)
// 		}
// 		// fmt.Println(PostIdFromHTML, "---------")
// 		var pos []post
// 		pos = displayPostsAndComments()

// 		AllLikes, AllDislikes := SumOfAllLikes(AllForumUsers())
// 		pos = DistLikesToPosts(pos, AllLikes, AllDislikes)
// 		allForumUnames := allForumUnames()
// 		var Chosen []post
// 		for i := 0; i < len(pos); i++ {
// 			if pos[i].PostID == PostIdFromHTML {
// 				Chosen = append(Chosen, pos[i])
// 			}
// 		}
// 		//********* IP ********
// 		duplicateIP = false
// 		if Chosen[0].IPs == "" {
// 			Chosen[0].IPs = GetOutboundIP().String()
// 			duplicateIP = true
// 		}
// 		if Chosen[0].IPs == GetOutboundIP().String() {
// 			duplicateIP = true
// 		}

// 		if !duplicateIP {
// 			Chosen[0].IPs += "-" + GetOutboundIP().String()
// 		}
// 		allIp := (strings.Split(Chosen[0].IPs, "-"))
// 		keys := make(map[string]bool)
// 		list := []string{}
// 		for _, entry := range allIp {
// 			if _, value := keys[entry]; !value {
// 				keys[entry] = true
// 				list = append(list, entry)
// 			}
// 		}
// 		Chosen[0].View = len(list)
// 		stmt, err := db.Prepare("UPDATE posts SET ips = ?	WHERE postID = ?;")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer stmt.Close()
// 		stmt.Exec(Chosen[0].IPs, Chosen[0].PostID)
// 		//********* IP ********

// 		urlPost = "postpage?postdetails=" + strID + "&postdetails=" + Chosen[0].Title
// 		fmt.Println(r.URL.Path)
// 		fmt.Println(urlPost)
// 		if r.URL.Path+"?postdetails="+strID+"&postdetails="+Chosen[0].Title != "/"+urlPost {
// 			http.Error(w, "404 Page Not Found", 404)
// 			return
// 		}
// 		Alllikes, Alldislikes := CommentSumOfAllLikes(AllForumUsers())

// 		Chosen[0].Comments = DistLikesToComments(Chosen[0].Comments, Alllikes, Alldislikes)
// 		if changingCom {
// 			userLikes := CountLikesByUser(curUser, "l")
// 			userDislikes := CountLikesByUser(curUser, "d")
// 			userComLikes := CommentCountLikesByUser(curUser, "l")
// 			userComDislikes := CommentCountLikesByUser(curUser, "d")
// 			for k := 0; k < len(userLikes); k++ {
// 				if Chosen[0].PostID == userLikes[k] {
// 					Chosen[0].LikedByCur = true
// 				}
// 			}
// 			for k := 0; k < len(userDislikes); k++ {
// 				if Chosen[0].PostID == userDislikes[k] {
// 					Chosen[0].DislikedByCur = true
// 				}
// 			}

// 			for i := 0; i < len(Chosen[0].Comments); i++ {
// 				for k := 0; k < len(userComLikes); k++ {
// 					if Chosen[0].Comments[i].CommentID == userComLikes[k] {
// 						Chosen[0].Comments[i].LikedByCur = true
// 					}
// 				}
// 			}
// 			for i := 0; i < len(Chosen[0].Comments); i++ {
// 				for k := 0; k < len(userComDislikes); k++ {
// 					if Chosen[0].Comments[i].CommentID == userComDislikes[k] {
// 						Chosen[0].Comments[i].DislikedByCur = true
// 					}
// 				}
// 			}
// 		}
// 		data := mainPageData{
// 			Posts:       Chosen,
// 			Userinfo:    curUser,
// 			ForumUnames: allForumUnames,
// 		}
// 		if changingCom {
// 			err = tpl.ExecuteTemplate(w, "post2.gohtml", data)
// 		} else {
// 			err = tpl.ExecuteTemplate(w, "post.gohtml", data)
// 		}

// 		if err != nil {
// 			http.Error(w, "Executing Error", http.StatusInternalServerError)
// 			return
// 		}
// 	} else if r.Method == "POST" {
// 		removeValue := r.FormValue("removePost")
// 		removeCom := r.FormValue("removeCom")
// 		editID := r.FormValue("editButton")
// 		editIDCom := r.FormValue("editButtonCom")
// 		fmt.Println(removeValue)
// 		if removeValue != "" {
// 			removeInt, err := strconv.Atoi(removeValue)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			pos := FindPostByPostId(removeInt)
// 			if pos.Author == curUser.Username {
// 				DeleteOnePost(removeInt)
// 				http.Redirect(w, r, "/", http.StatusSeeOther)
// 			} else {
// 				http.Error(w, "You are not authorized", 400)
// 			}

// 		}
// 		if removeCom != "" {
// 			removeInt, err := strconv.Atoi(removeCom)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			com := FindComByComId(removeInt)
// 			if com.Author == curUser.Username {
// 				DeleteOneCom(removeInt)
// 				http.Redirect(w, r, "/", http.StatusSeeOther)
// 			} else {
// 				http.Error(w, "You are not authorized", 400)
// 			}

// 		}
// 		if editID != "" {

// 			var newPost post
// 			posID, err := strconv.Atoi(editID)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			pos := FindPostByPostId(posID)

// 			postCon := r.FormValue("postContentE")
// 			postTitle := r.FormValue("postTitleE")
// 			newPost = pos
// 			if postCon != "" {
// 				newPost.Content = postCon
// 			}
// 			if postTitle != "" {
// 				newPost.Title = postTitle
// 			}
// 			newUrl := "localhost:8080/postpage?postdetails=" + editID + "&postdetails=" + newPost.Title
// 			newPost.PostTime = time.Now()
// 			newPost.URL = newUrl
// 			if pos.Author == curUser.Username {
// 				EditPost(newPost)
// 				http.Redirect(w, r, "/", http.StatusSeeOther)
// 			} else {
// 				http.Error(w, "You are not authorized", 400)
// 			}

// 		}
// 		if editIDCom != "" {
// 			var newCom comment
// 			comID, err := strconv.Atoi(editIDCom)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			com := FindComByComId(comID)

// 			comCon := r.FormValue("comContentE")

// 			newCom = com
// 			if comCon != "" {
// 				newCom.Content = comCon
// 			}

// 			newCom.CommentTime = time.Now()
// 			if com.Author == curUser.Username {
// 				EditCom(newCom)
// 				http.Redirect(w, r, "/", http.StatusSeeOther)
// 			} else {
// 				http.Error(w, "You are not authorized", 400)
// 			}
// 		}
// 		processPost(r, curUser)
// 		processComment(r, curUser)
// 		http.Redirect(w, r, urlPost, http.StatusSeeOther)

// 	} else {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 	}
// }

// func CategoryPageHandler(w http.ResponseWriter, r *http.Request) {
// 	curUser := obtainCurUserFormCookie(r)
// 	if r.Method == "GET" {
// 		w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 		tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/header2.gohtml", "./templates/footer.gohtml", "./templates/categories.gohtml")
// 		if err != nil {
// 			http.Error(w, "Parsing Error", http.StatusInternalServerError)
// 			return
// 		}

// 		var pos []post
// 		category := r.FormValue("categoryAllPosts")
// 		pos = filCatDisplayPostsAndComments(category)

// 		allForumUnames := allForumUnames()
// 		data := mainPageData{
// 			Posts:       pos,
// 			Userinfo:    curUser,
// 			ForumUnames: allForumUnames,
// 		}
// 		// fmt.Println("---------", forumUser)
// 		err = tpl.ExecuteTemplate(w, "categories.gohtml", data)
// 		if err != nil {
// 			http.Error(w, "Executing Error", http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 	}
// }

// func NotiPageHandler(w http.ResponseWriter, r *http.Request) {
// 	CurUser := obtainCurUserFormCookie(r)
// 	users := AllForumUsers()
// 	for i := 0; i < len(users); i++ {
// 		if users[i].Username == CurUser.Username {
// 			CurUser.Notifymsg = users[i].Notifymsg
// 			CurUser.Notifyview = users[i].Notifyview
// 		}
// 	}
// 	var NewCodes []string
// 	if r.Method == "GET" {
// 		w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 		tpl, err := template.ParseFiles("./templates/header2.gohtml", "./templates/footer.gohtml", "./templates/notif.gohtml")
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Parsing Error", http.StatusInternalServerError)
// 			return
// 		}
// 		CurUser, NewCodes = UpdateNotif(CurUser)
// 		fmt.Println("************NOTIFICATION", CurUser.NotifMessageShow)
// 		NewCodesStr := strings.Join(NewCodes, "#")
// 		CurUser.Notifyview += "#" + NewCodesStr
// 		stmt, err := db.Prepare("UPDATE users SET Notifyview = ?	WHERE username = ?;")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer stmt.Close()
// 		stmt.Exec(CurUser.Notifyview, CurUser.Username)

// 		// fmt.Println("---------", forumUser)
// 		err = tpl.ExecuteTemplate(w, "notif.gohtml", CurUser)
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Executing Error", http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 	}
// }

// // func DeleteHandler(w http.ResponseWriter, r *http.Request) {
// // 	// for testing purpose
// // 	if r.Method == http.MethodGet {
// // 		tpl, err := template.ParseFiles("./templates/delete.gohtml", "./templates/footer.gohtml", "./templates/header.gohtml")
// // 		if err != nil {
// // 			log.Fatal(err)
// // 		}
// // 		tpl.ExecuteTemplate(w, "delete.gohtml", nil)
// // 	}
// // 	if r.Method == http.MethodPost {
// // 		deleteUser(r)
// // 	}
// // }

// // func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
// // 	tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/footer.gohtml", "./templates/notFound.gohtml")
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	tpl.ExecuteTemplate(w, "notFound.gohtml", nil)
// // }
// func ActivityPageHandler(w http.ResponseWriter, r *http.Request) {
// 	var act Activity
// 	CurUser := obtainCurUserFormCookie(r)
// 	users := AllForumUsers()
// 	for i := 0; i < len(users); i++ {
// 		if users[i].Username == CurUser.Username {
// 			CurUser = users[i]
// 		}
// 	}
// 	// act = FillActivity(CurUser)
// 	LikedInt := CountLikesByUser(CurUser, "l")
// 	DlikedInt := CountLikesByUser(CurUser, "d")
// 	ComLikedInt := CommentCountLikesByUser(CurUser, "l")
// 	fmt.Println(ComLikedInt, "ComLikedInt")
// 	ComDlikedInt := CommentCountLikesByUser(CurUser, "d")
// 	act.Username = CurUser.Username
// 	act = ActFindingPostAndCom(LikedInt, act, "Post", "Liked")
// 	act = ActFindingPostAndCom(DlikedInt, act, "Post", "Disliked")
// 	act = ActFindingPostAndCom(ComLikedInt, act, "Com", "Liked")
// 	act = ActFindingPostAndCom(ComDlikedInt, act, "Com", "Disliked")
// 	act = CreatedPostandCom(act)
// 	fmt.Println(act.LikedCom)
// 	if r.Method == "GET" {
// 		w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 		tpl, err := template.ParseFiles("./templates/header2.gohtml", "./templates/footer.gohtml", "./templates/activity.gohtml")
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Parsing Error", http.StatusInternalServerError)
// 			return
// 		}
// 		err = tpl.ExecuteTemplate(w, "activity.gohtml", act)
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Executing Error", http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 	}
// }

// func EditorRemovePageHandler(w http.ResponseWriter, r *http.Request) {
// 	CurUser := obtainCurUserFormCookie(r)
// 	users := AllForumUsers()
// 	for i := 0; i < len(users); i++ {
// 		if users[i].Username == CurUser.Username {
// 			CurUser = users[i]
// 		}
// 	}
// }
