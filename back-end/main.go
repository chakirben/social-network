package main

import (
	"fmt"
	"net/http"

	auth "socialN/Handlers/auth"
	db "socialN/dataBase"
	Comment"socialN/Handlers/comments"
)

func setupHandlers() {

	//Serve Images
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	// auth
	http.HandleFunc("/api/checkAuth", AccessMiddleware(auth.CheckAuth))
	http.HandleFunc("/api/register",AccessMiddleware(auth.RegisterUser))
	http.HandleFunc("/api/login", AccessMiddleware(auth.LogUser))
	http.HandleFunc("/api/logout", AccessMiddleware(auth.LogoutHandler))
	http.HandleFunc("/api/profile", AccessMiddleware(auth.ProfileHandler))

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

	//comments
	http.HandleFunc("/api/GetComments", AccessMiddleware(Comment.GetCommentsHandler))
	http.HandleFunc("/api/SetComment", AccessMiddleware(Comment.SetCommentHandler))
	// http.HandleFunc("/api/Like", handlers.ReactionHandler)
	http.HandleFunc("/api/Profile", auth.ProfileHandler)
	http.HandleFunc("/api/CheckAuth", auth.CheckAuth)
}

func main() {
	db.DbInit()
	setupHandlers()
	fmt.Println("https://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func Shandler(fun http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	auth.ValidateSession(r, db.SocialDB)
	fun(w, r)
}

func SessionMiddleware(fun http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.ValidateSession(r, db.SocialDB)
		if err != nil {
			fmt.Println("uno")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fun(w, r)
	}
}
func AccessMiddleware(fun http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		fun(w, r)
	}
}
