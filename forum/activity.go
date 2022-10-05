package forum

func ActFindingPostAndCom(Likes []int, act Activity, str, str2 string) Activity {
	if str == "Post" {
		var likedPosts []post
		PostCom := displayPostsAndComments()
		for i := 0; i < len(PostCom); i++ {
			for k := 0; k < len(Likes); k++ {
				if PostCom[i].PostID == Likes[k] {
					likedPosts = append(likedPosts, PostCom[i])
				}
			}
		}
		if str2 == "Liked" {
			act.LikedPost = likedPosts
		} else {
			act.DlikedPost = likedPosts
		}
		return act
	} else {
		var likedPosts []comment
		PostCom := displayComs()
		for i := 0; i < len(PostCom); i++ {
			for k := 0; k < len(Likes); k++ {
				if PostCom[i].CommentID == Likes[k] {
					likedPosts = append(likedPosts, PostCom[i])
				}
			}
		}
		if str2 == "Liked" {
			act.LikedCom = likedPosts
		} else {
			act.DlikedCom = likedPosts
		}
		return act
	}
}

func CreatedPostandCom(act Activity) Activity {
	po := displayPostsAndComments()
	com := displayComs()
	for i := 0; i < len(po); i++ {
		if po[i].Author == act.Username {
			act.Post = append(act.Post, po[i])
		}
	}
	for i := 0; i < len(com); i++ {
		if com[i].Author == act.Username {
			act.Com = append(act.Com, com[i])
		}
	}
	return act
}
