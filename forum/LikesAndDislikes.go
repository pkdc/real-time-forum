package forum

import (
	"fmt"
	"strconv"
	"strings"
)

func CountLikesByUser(curUser user, str string) []int {
	var blankslice []string
	var curUserLikesInt []int
	var curUserLikes2 []string
	var curUserLikes []string

	if str == "l" {
		curUserLikes = strings.Split(curUser.LikedPost, "-")
	} else if str == "d" {
		curUserLikes = strings.Split(curUser.DislikedPost, "-")
	}
	curUserLikes = append(blankslice, curUserLikes[+1:]...)
	for k, v := range dup_count(curUserLikes) {
		if v%2 != 0 && v != 0 {
			curUserLikes2 = append(curUserLikes2, k)
		}
	}
	for _, i := range curUserLikes2 {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		curUserLikesInt = append(curUserLikesInt, j)
	}

	return curUserLikesInt
}

func SumOfAllLikes(users []user) ([]int, []int) {
	var AllLikes []int
	var AllDislikes []int
	for i := 0; i < len(users); i++ {
		AllLikes = append(AllLikes, CountLikesByUser(users[i], "l")...)
	}
	for i := 0; i < len(users); i++ {
		AllDislikes = append(AllDislikes, CountLikesByUser(users[i], "d")...)
	}
	return AllLikes, AllDislikes
}

func DistLikesToPosts(pos []post, Alllikes, AllDislikes []int) []post {
	fmt.Println("***************used")
	curUserLikes := make([]string, len(Alllikes))
	curUserDislikes := make([]string, len(AllDislikes))
	for i, t := range Alllikes {
		curUserLikes[i] = strconv.Itoa(t)
	}
	for i, t := range AllDislikes {
		curUserDislikes[i] = strconv.Itoa(t)
	}
	for i := 0; i < len(pos); i++ {
		for k, v := range dup_count(curUserLikes) {
			if k != "" {
				d, err := strconv.Atoi(k)
				if err != nil {
					continue
				}
				fmt.Printf("whichPost %s, howmanytimes liked %d", k, v)
				if pos[i].PostID == d {
					pos[i].Likes = v
				}
			} else {
				pos[i].Likes = 0
			}
		}
	}
	for i := 0; i < len(pos); i++ {
		for k, v := range dup_count(curUserDislikes) {
			if k != "" {
				d, err := strconv.Atoi(k)
				if err != nil {
					continue
				}
				fmt.Printf("whichPost %s, howmanytimes disliked %d", k, v)
				if pos[i].PostID == d {
					pos[i].Dislikes = v
				}
			} else {
				pos[i].Dislikes = 0
			}
		}
	}
	return pos
}

func CommentCountLikesByUser(curUser user, str string) []int {
	var blankslice []string
	var curUserLikesInt []int
	var curUserLikes2 []string
	var curUserLikes []string
	if str == "l" {
		curUserLikes = strings.Split(curUser.LikedComments2, "-")
	} else if str == "d" {
		curUserLikes = strings.Split(curUser.DislikedComments2, "-")
	}
	curUserLikes = append(blankslice, curUserLikes[+1:]...)
	for k, v := range dup_count(curUserLikes) {
		if v%2 != 0 && v != 0 {
			curUserLikes2 = append(curUserLikes2, k)
		}
	}
	for _, i := range curUserLikes2 {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		curUserLikesInt = append(curUserLikesInt, j)
	}
	return curUserLikesInt
}

func CommentSumOfAllLikes(users []user) ([]int, []int) {
	var AllLikes []int
	var AllDislikes []int
	for i := 0; i < len(users); i++ {
		AllLikes = append(AllLikes, CommentCountLikesByUser(users[i], "l")...)
	}
	for i := 0; i < len(users); i++ {
		AllDislikes = append(AllDislikes, CommentCountLikesByUser(users[i], "d")...)
	}
	return AllLikes, AllDislikes
}

func DistLikesToComments(com []comment, Alllikes, AllDislikes []int) []comment {
	curUserLikes := make([]string, len(Alllikes))
	curUserDislikes := make([]string, len(AllDislikes))
	for i, t := range Alllikes {
		curUserLikes[i] = strconv.Itoa(t)
	}
	for i, t := range AllDislikes {
		curUserDislikes[i] = strconv.Itoa(t)
	}
	for i := 0; i < len(com); i++ {
		for k, v := range dup_count(curUserLikes) {
			if k != "" {
				d, err := strconv.Atoi(k)
				if err != nil {
					continue
				}
				fmt.Printf("whichCOM %s, howmanytimes liked %d", k, v)
				if com[i].CommentID == d {
					com[i].Likes = v
				}
			} else {
				com[i].Likes = 0
			}
		}
	}
	for i := 0; i < len(com); i++ {
		fmt.Println(dup_count(curUserDislikes))
		for k, v := range dup_count(curUserDislikes) {
			if k != "" {
				d, err := strconv.Atoi(k)
				if err != nil {
					continue
				}
				fmt.Printf("whichCOM %s, howmanytimes disliked %d", k, v)
				if com[i].CommentID == d {
					com[i].Dislikes = v
				}
			} else {
				com[i].Dislikes = 0
			}
		}
	}
	return com
}

func CheckLikesAndDislikes(usr user, posID int, str string) string {
	var likesordislikes string
	var likes []int
	if str == "l" {
		likes = CountLikesByUser(usr, str)
		likesordislikes = usr.LikedPost
	} else if str == "d" {
		likes = CountLikesByUser(usr, str)
		likesordislikes = usr.DislikedPost
	} else if str == "ComL" {
		likes = CommentCountLikesByUser(usr, "l")
		likesordislikes = usr.LikedComments2
	} else if str == "ComD" {
		likes = CommentCountLikesByUser(usr, "d")
		likesordislikes = usr.DislikedComments2
	}
	posStr := strconv.Itoa(posID)

	for k := 0; k < len(likes); k++ {
		if posID == likes[k] {
			likedPost := strings.Split(likesordislikes, "")
			for i := 0; i < len(likesordislikes); i++ {
				if likedPost[i] == posStr && i != 0 {
					likedPost = append(likedPost[:i-1], likedPost[i+1:]...)
					break
				} else if likedPost[i] == posStr && i == 0 {
					likedPost = append(likedPost[:i], likedPost[i+2:]...)
					break
				}
			}
			likesordislikes = strings.Join(likedPost, "")
		}
	}

	return likesordislikes
}
