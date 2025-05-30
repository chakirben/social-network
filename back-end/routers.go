package main

import (
	"fmt"
	"net/http"

	"socialN/Handlers/auth"
	"socialN/Handlers/followers"
	"socialN/Handlers/notification"
	"socialN/Handlers/ws"

	h "socialN/Handlers"
	Comment "socialN/Handlers/comments"
	event "socialN/Handlers/events"
	Group "socialN/Handlers/groups"
	Post "socialN/Handlers/posts"
	profile "socialN/Handlers/profile"

	Search "socialN/Handlers/Search"

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
	http.HandleFunc("/api/SetComment", AccessMiddleware(Comment.SetCommentHandler))

	// posts
	http.HandleFunc("/api/GetCreatedPosts", AccessMiddleware(Post.GetCreatedPostsHandler))
	http.HandleFunc("/api/GetOnePost", AccessMiddleware(SessionMiddleware(Post.GetPostHandler)))
	// http.HandleFunc("/api/GetOnePost", AccessMiddleware(SessionMiddleware(Post.CreatePostHandler)))
	http.HandleFunc("/api/CreatePost", AccessMiddleware(SessionMiddleware(Post.CreatePostHandler)))
	http.HandleFunc("/api/GetPosts", AccessMiddleware(SessionMiddleware(Post.GetPostsHandler)))
	// http.HandleFunc("/api/GetLikedPosts", Post.GetLikedPostsHandler)

	// Events
	http.HandleFunc("/api/CreateEvent", AccessMiddleware(event.SetEventHandler))
	http.HandleFunc("/api/SetAttendance", AccessMiddleware(event.SetAttendanceHandler))
	http.HandleFunc("/api/GetHomeEvents", AccessMiddleware(event.GetHomeEventsHandler))
	http.HandleFunc("/api/GetGroupEvents", AccessMiddleware(event.GetGroupEventsHandler))

	// groups Creat_Groups
	http.HandleFunc("/api/CreatGroup", AccessMiddleware(SessionMiddleware(Group.Creat_Groups)))
	http.HandleFunc("/api/JoinGroup", AccessMiddleware(SessionMiddleware(Group.JoinGroup)))
	http.HandleFunc("/api/MyGroups", AccessMiddleware(SessionMiddleware(Group.GetMyGroups)))
	http.HandleFunc("/api/NotMyGroups", AccessMiddleware(SessionMiddleware(Group.GetGroupsUserNotJoined)))
	http.HandleFunc("/api/PostsGroups", AccessMiddleware(SessionMiddleware(Group.GetPostGroups)))
	http.HandleFunc("/api/RequestToJoinGroups", AccessMiddleware(SessionMiddleware(Group.Req_To_Join_Groups)))

	// SearchData
	http.HandleFunc("/api/SearchData", AccessMiddleware(SessionMiddleware(Search.SearchData)))

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
	//	http.HandleFunc("/api/follow", AccessMiddleware(u.Follow))
	http.HandleFunc("/api/follow", AccessMiddleware(followers.HandleFollow))
	http.HandleFunc("/api/acceptFollowRequest", AccessMiddleware(followers.AcceptFollowRequest))
	http.HandleFunc("/api/declineFollowRequest", AccessMiddleware(followers.DeclineFollowRequest))

	// profile
	http.HandleFunc("/api/profile", AccessMiddleware(profile.GetData))

	// reactions
	http.HandleFunc("/api/reaction", AccessMiddleware(h.ReactionHandler))

	// users
	http.HandleFunc("/api/getUserData", AccessMiddleware(u.GetCurrentUserData))
	http.HandleFunc("/api/getUnfollowedUsers", AccessMiddleware(u.GetUnfollowedUsers))

	http.HandleFunc("/api/updatePrivacy", AccessMiddleware(u.SetPrivacy))
	http.HandleFunc("/api/getFollowersList", AccessMiddleware(u.GetFollowersListHandler))

	// notifications
	http.HandleFunc("/api/getNotifications", AccessMiddleware(notification.GetNotifications))
	http.HandleFunc("/api/ws", ws.OpenWsConn)
	http.HandleFunc("/api/online" , AccessMiddleware(ws.GetOnlineUsers))
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
