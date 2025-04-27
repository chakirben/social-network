package main

import (
	"fmt"
	"net/http"

	auth "socialN/Handlers/auth"
	db "socialN/dataBase"
)

func setupHandlers() {
	// auth
	http.HandleFunc("/api/checkAuth", auth.CheckAuth)
	http.HandleFunc("/api/register",auth.RegisterUser)
	http.HandleFunc("/api/Login", SessionMiddleware(auth.LogUser))
	http.HandleFunc("/api/Logout", SessionMiddleware(auth.LogoutHandler))

	// // chat
	// http.HandleFunc("/api/Chat", chat.ChatHandler)
	// http.HandleFunc("/api/GetMessages", chat.GetMessagesHandler)
	// http.HandleFunc("/api/GetDiscussions", chat.GetDiscussionsListHandler)
	// http.HandleFunc("/api/GetOnlineUsers", chat.GetOnlineUsersHandler)

	// // posts
	// http.HandleFunc("/api/CreatePost", Post.SetPostHandler)
	// http.HandleFunc("/api/GetPosts", Post.GetPostsHandler)
	// http.HandleFunc("/api/GetPost", Post.GetPostHandler)
	// http.HandleFunc("/api/GetLikedPosts", Post.GetLikedPostsHandler)
	// http.HandleFunc("/api/GetCreatedPosts", Post.GetCreatedPostsHandler)

	// // comments
	// http.HandleFunc("/api/GetComments", Comment.GetCommentsHandler)
	// http.HandleFunc("/api/SetComment", Comment.SetCommentHandler)
	// http.HandleFunc("/api/Like", handlers.ReactionHandler)
	// http.HandleFunc("/api/Profile", auth.ProfileHandler)
	// http.HandleFunc("/api/CheckAuth", auth.CheckAuth)
}

func main() {
	db.DbInit()
	setupHandlers()
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Println("in here")})
	fmt.Println("https://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func Shandler(fun http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	auth.ValidateSession(r, db.SocialDB)
	fun(w, r)
}

func SessionMiddleware(fun http.HandlerFunc) http.HandlerFunc {
	fmt.Println("sess")
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hhhh")
		_, err := auth.ValidateSession(r, db.SocialDB)
		if err != nil {
			fmt.Println("uno")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fun(w, r)
	}
}
