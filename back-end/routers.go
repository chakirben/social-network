package main

import (
	"fmt"
	"net/http"

	"socialN/Handlers/auth"
	"socialN/Handlers/followers"

	Comment "socialN/Handlers/comments"
	event "socialN/Handlers/events"
	Group "socialN/Handlers/groups"
	Post "socialN/Handlers/posts"
	profile "socialN/Handlers/profile"
	u "socialN/Handlers/users"
	db "socialN/dataBase"
)

func SetupHandlers() {
	// Serve Images
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// auth
	http.HandleFunc("/api/checkAuth", AccessMiddleware(auth.CheckAuth))
	http.HandleFunc("/api/register", AccessMiddleware(auth.RegisterUser))
	http.HandleFunc("/api/login", AccessMiddleware(auth.LogUser))
	http.HandleFunc("/api/logout", AccessMiddleware(auth.LogoutHandler))
	// http.HandleFunc("/api/profile", AccessMiddleware(auth.ProfileHandler))

	// comments
	http.HandleFunc("/api/GetComments", AccessMiddleware(Comment.GetCommentsHandler))
	http.HandleFunc("/api/SetComment", AccessMiddleware(SessionMiddleware(Comment.SetCommentHandler)))

	// posts
	http.HandleFunc("/api/GetCreatedPosts", AccessMiddleware(Post.GetCreatedPostsHandler))
	http.HandleFunc("/api/GetOnePost", AccessMiddleware(SessionMiddleware(Post.GetPostHandler)))
	// http.HandleFunc("/api/GetOnePost", AccessMiddleware(SessionMiddleware(Post.CreatePostHandler)))
	http.HandleFunc("/api/CreatePost", AccessMiddleware(Post.CreatePostHandler))
	http.HandleFunc("/api/GetPosts", AccessMiddleware(Post.GetPostsHandler))
	// http.HandleFunc("/api/GetLikedPosts", Post.GetLikedPostsHandler)

	// Events
	http.HandleFunc("/api/CreateEvent", event.SetEventHandler)
	http.HandleFunc("/api/SetAttendance", event.SetAttendanceHandler)
	http.HandleFunc("/api/GetEvents", event.GetEventsHandler)

	// groups Creat_Groups
	http.HandleFunc("/api/CreatGroup", Group.Creat_Groups)
	http.HandleFunc("/api/JoinGroup", AccessMiddleware(SessionMiddleware(Group.JoinGroup)))
	http.HandleFunc("/api/MyGroups", AccessMiddleware(Group.GetMyGroups))
	http.HandleFunc("/api/NotMyGroups", AccessMiddleware(Group.GetGroupsUserNotJoined))
	http.HandleFunc("/api/RequestToJoinGroups", Group.Req_To_Join_Groups)

	// // chat
	// http.HandleFunc("/api/Chat", chat.ChatHandler)
	// http.HandleFunc("/api/GetMessages", chat.GetMessagesHandler)
	// http.HandleFunc("/api/GetDiscussions", chat.GetDiscussionsListHandler)
	// http.HandleFunc("/api/GetOnlineUsers", chat.GetOnlineUsersHandler)

	// http.HandleFunc("/api/Like", handlers.ReactionHandler)
	http.HandleFunc("/api/Profile", AccessMiddleware(auth.ProfileHandler))
	// http.HandleFunc("/api/CheckAuth", AccessMiddleware(auth.CheckAuth))
	// http.HandleFunc("/api/Profile", auth.ProfileHandler)
	http.HandleFunc("/api/CheckAuth", auth.CheckAuth)

	// this just for testing you can delete it
	// the function has the id of loggedin user as a parameter, you can get it from session
	fmt.Print("Followers of loggedin user")
	fmt.Println(followers.GetFollowedUsers(2))

	// follows
	http.HandleFunc("/api/follow", followers.HandleFollow)
	http.HandleFunc("/api/acceptFollowRequest", followers.AcceptFollowRequest)

	// profile
	http.HandleFunc("/api/profile", profile.GetData)

	//
	http.HandleFunc("/api/getUserData", AccessMiddleware(u.GetCurrentUserData))
	http.HandleFunc("/api/updatePrivacy", AccessMiddleware(u.SetPrivacy))
	http.HandleFunc("/api/getFollowersList", AccessMiddleware(u.GetFollowersListHandler))
}

func SessionMiddleware(fun http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.ValidateSession(r, db.SocialDB)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fun(w, r)
	}
}

func AccessMiddleware(fun http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", `http://localhost:3000`)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		fun(w, r)
	}
}
