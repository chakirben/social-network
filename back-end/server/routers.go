package main

import (
	"net/http"

	"socialN/app/auth"
	"socialN/app/followers"
	"socialN/app/notification"
	"socialN/app/ws"

	h "socialN/app"
	Comment "socialN/app/comments"
	event "socialN/app/events"
	Group "socialN/app/groups"
	Post "socialN/app/posts"
	profile "socialN/app/profile"

	Search "socialN/app/Search"

	u "socialN/app/users"

	db "socialN/dataBase"
)

func SetupHandlers() {
	// Serve Images
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// auth
	http.HandleFunc("/api/checkAuth", (auth.CheckAuth))
	http.HandleFunc("/api/register", (auth.RegisterUser))
	http.HandleFunc("/api/login", (auth.LogUser))
	http.HandleFunc("/api/logout", (auth.LogoutHandler))
	// http.HandleFunc("/api/profile", (auth.ProfileHandler))

	// comments
	http.HandleFunc("/api/GetComments", (Comment.GetCommentsHandler))
	http.HandleFunc("/api/SetComment", (Comment.SetCommentHandler))

	// posts
	http.HandleFunc("/api/GetCreatedPosts", (Post.GetCreatedPostsHandler))
	http.HandleFunc("/api/GetOnePost", (SessionMiddleware(Post.GetPostHandler)))
	// http.HandleFunc("/api/GetOnePost", (SessionMiddleware(Post.CreatePostHandler)))
	http.HandleFunc("/api/CreatePost", (SessionMiddleware(Post.CreatePostHandler)))
	http.HandleFunc("/api/GetPosts", (SessionMiddleware(Post.GetPostsHandler)))
	// http.HandleFunc("/api/GetLikedPosts", Post.GetLikedPostsHandler)

	// Events
	http.HandleFunc("/api/CreateEvent", (event.SetEventHandler))
	http.HandleFunc("/api/SetAttendance", (event.SetAttendanceHandler))
	http.HandleFunc("/api/GetHomeEvents", (event.GetHomeEventsHandler))
	http.HandleFunc("/api/GetGroupEvents", (event.GetGroupEventsHandler))

	// groups Creat_Groups
	http.HandleFunc("/api/CreatGroup", (SessionMiddleware(Group.Creat_Groups)))
	http.HandleFunc("/api/JoinGroup", (SessionMiddleware(Group.JoinGroup)))
	http.HandleFunc("/api/MyGroups", (SessionMiddleware(Group.GetMyGroups)))
	http.HandleFunc("/api/NotMyGroups", (SessionMiddleware(Group.GetGroupsUserNotJoined)))
	http.HandleFunc("/api/PostsGroups", (SessionMiddleware(Group.GetPostGroups)))
	http.HandleFunc("/api/RequestToJoinGroups", (SessionMiddleware(Group.Req_To_Join_Groups)))
	http.HandleFunc("/api/CancelRequestToJoinGroups", (SessionMiddleware(Group.CancelRequestToJoinGroups)))
	http.HandleFunc("/api/CancelInviteToGroups", (SessionMiddleware(Group.CancelInviteToGroups)))
	http.HandleFunc("/api/getFollowers", (Group.GetFollowersList))
	http.HandleFunc("/api/InfiteTheFollowers", (SessionMiddleware(Group.InviteTheFollowers)))

	// SearchData
	http.HandleFunc("/api/SearchData", (SessionMiddleware(Search.SearchData)))

	// // chat
	// http.HandleFunc("/api/Chat", chat.ChatHandler)
	// http.HandleFunc("/api/GetMessages", chat.GetMessagesHandler)
	// http.HandleFunc("/api/GetDiscussions", chat.GetDiscussionsListHandler)
	// http.HandleFunc("/api/GetOnlineUsers", chat.GetOnlineUsersHandler)

	// http.HandleFunc("/api/Like", app.ReactionHandler)
	http.HandleFunc("/api/Profile", (auth.ProfileHandler))
	// http.HandleFunc("/api/CheckAuth", (auth.CheckAuth))
	// http.HandleFunc("/api/Profile", auth.ProfileHandler)
	http.HandleFunc("/api/CheckAuth", auth.CheckAuth)

	// follows
	http.HandleFunc("/api/follow", (followers.FollowHandler))
	http.HandleFunc("/api/acceptFollowRequest", (followers.AcceptFollowRequest))
	http.HandleFunc("/api/declineFollowRequest", (followers.DeclineFollowRequest))

	// profile
	http.HandleFunc("/api/profile", (profile.GetData))

	// reactions
	http.HandleFunc("/api/reaction", (h.ReactionHandler))

	// users
	http.HandleFunc("/api/getUserData", (u.GetCurrentUserData))
	http.HandleFunc("/api/getUnfollowedUsers", (u.GetUnfollowedUsers))
	http.HandleFunc("/api/updatePrivacy", (u.SetPrivacy))
	http.HandleFunc("/api/getFollowersList", (u.GetFollowersListHandler))
	http.HandleFunc("/api/followersList", (followers.GetListHandler))

	// notifications
	http.HandleFunc("/api/getNotifications", (notification.GetNotifications))
	http.HandleFunc("/api/respondToNotification", (notification.RespondtoNotification))

	// chat
	http.HandleFunc("/api/ws", ws.OpenWsConn)
	http.HandleFunc("/api/GetDiscussionList", (ws.GetAllDiscussionsHandler))
	http.HandleFunc("/api/online", (ws.GetOnlineUsers))
	http.HandleFunc("/api/fetchMessages", (ws.GetMessagesHandler))
	http.HandleFunc("/api/friendsAndGroups", (ws.GetUserConnectionsHandler))
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
