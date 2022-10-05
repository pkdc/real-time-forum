package forum

import (
	"html/template"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func findAuthor(posID int) (string, user) {
	var SelectedUser user
	var authorName string
	var msg string
	posIDstr := strconv.Itoa(posID)
	po := displayPostsAndComments()
	usr := AllForumUsers()
	for i := 0; i < len(po); i++ {
		if po[i].PostID == posID {
			authorName = po[i].Author
			msg = "#" + "localhost:8080/postpage?postdetails=" + posIDstr + "&postdetails=" + po[i].Title + "#"
		}
	}
	for i := 0; i < len(usr); i++ {
		if usr[i].Username == authorName {
			SelectedUser = usr[i]
		}
	}

	return msg, SelectedUser
}

func findCommentAuthor(comID int) (string, user) {
	var SelectedUser user
	var authorName string
	var msg string
	var posID int
	po := displayPostsAndComments()
	for i := 0; i < len(po); i++ {
		for k := 0; k < len(po[i].Comments); k++ {
			if po[i].Comments[k].CommentID == comID {
				posID = po[i].PostID
				authorName = po[i].Comments[k].Author
			}
		}
	}
	posIDstr := strconv.Itoa(posID)
	usr := AllForumUsers()
	for i := 0; i < len(po); i++ {
		if po[i].PostID == posID {
			msg = "#" + "localhost:8080/postpage?postdetails=" + posIDstr + "&postdetails=" + po[i].Title + "#"
		}
	}
	for i := 0; i < len(usr); i++ {
		if usr[i].Username == authorName {
			SelectedUser = usr[i]
		}
	}

	return msg, SelectedUser
}

// func showNotifications(usr user) user {
// 	msg := usr.NotifMessage
// 	view := usr.NotifView
// 	msg2 := strings.Split(msg, "#")
// 	view2 := strings.Split(view, "#")
// 	for i := 0; i < len(msg2); i++ {
// 		for k := 0; k < len(view2); k++ {
// 			if msg2[i] == view2[k] {
// 				msg2[i-2] = ""
// 				msg2[i-1] = ""
// 			}
// 		}
// 		usr.NotifMessage = strings.Join(msg2, "#")
// 	}
// 	return usr
// }

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func UpdateNotif(curUser user) (user, []string) {
	var Viewcodes []string
	var SeenCodes []int
	var NewMsg []string
	var NewCodes []string
	var NotifMessageShow []string
	var NotifMessageLink []string
	msg := curUser.Notifymsg
	msgSlc := strings.Split(msg, "#")
	for i := 2; i < len(msgSlc); i += 3 {
		Viewcodes = append(Viewcodes, msgSlc[i])
	}
	curUsrCodes := curUser.Notifyview
	curUsrCodesSlc := strings.Split(curUsrCodes, "#")
	for i := 0; i < len(Viewcodes); i++ {
		for k := 0; k < len(curUsrCodesSlc); k++ {
			if Viewcodes[i] == curUsrCodesSlc[k] {
				SeenCodes = append(SeenCodes, 3*i, (3*i)+1, (3*i)+2)
			}
		}
	}
	NewMsg = remove(msgSlc, SeenCodes)
	NewMsg = NewMsg[:len(NewMsg)-1]

	for i := 0; i < len(NewMsg); i += 3 {
		NotifMessageShow = append(NotifMessageShow, NewMsg[i])
	}
	for i := 1; i < len(NewMsg); i += 3 {
		NotifMessageLink = append(NotifMessageLink, NewMsg[i])
	}
	for i := 2; i < len(NewMsg); i += 3 {
		NewCodes = append(NewCodes, NewMsg[i])
	}

	curUser.NotifMessageShow = SafeUrl(NotifMessageShow, NotifMessageLink)
	return curUser, NewCodes
}

func remove(slice []string, s []int) []string {
	k := 0
	for i := 0; i < len(s); i++ {
		slice = append(slice[:s[i]-k], slice[s[i]-k+1:]...)
		k++
	}
	return slice
}

func SafeUrl(msg, link []string) map[string]template.URL {
	ShowMap := make(map[string]template.URL, len(msg)*2)
	var Slc []template.URL
	for i := 0; i < len(link); i++ {
		Slc = append(Slc, template.URL(link[i]))
	}

	for k := 0; k < len(Slc); k++ {
		intK := strconv.Itoa(k + 1)
		ShowMap[intK+"-"+msg[k]] = Slc[k]

	}

	return ShowMap
}
