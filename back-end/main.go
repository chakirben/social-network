package main

import (
	"fmt"
	"net/http"

	auth "socialN/Handlers/auth"
	Comment "socialN/Handlers/comments"
	Post "socialN/Handlers/posts"
	db "socialN/dataBase"
	event "socialN/Handlers/events"
)

func setupHandlers() {
	// Serve Images
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// auth
	http.HandleFunc("/api/checkAuth", AccessMiddleware(auth.CheckAuth))
	http.HandleFunc("/api/register", AccessMiddleware(auth.RegisterUser))
	http.HandleFunc("/api/login", AccessMiddleware(auth.LogUser))
	http.HandleFunc("/api/logout", AccessMiddleware(auth.LogoutHandler))
	http.HandleFunc("/api/profile", AccessMiddleware(auth.ProfileHandler))

	// comments
	http.HandleFunc("/api/GetComments", AccessMiddleware(SessionMiddleware(Comment.GetCommentsHandler)))
	http.HandleFunc("/api/SetComment", AccessMiddleware(SessionMiddleware(Comment.SetCommentHandler)))

	// posts
	http.HandleFunc("/api/GetCreatedPosts", AccessMiddleware(SessionMiddleware(Post.GetCreatedPostsHandler)))
	http.HandleFunc("/api/GetOnePost", AccessMiddleware(SessionMiddleware(Post.GetPostHandler)))
	// http.HandleFunc("/api/CreatePost", Post.SetPostHandler)
	http.HandleFunc("/api/GetPosts", Post.GetPostsHandler)
	// http.HandleFunc("/api/GetLikedPosts", Post.GetLikedPostsHandler)

	// Events
	http.HandleFunc("/api/CreateEvent", event.SetEventHandler)
	http.HandleFunc("/api/SetAttendance", event.SetAttendanceHandler)

	
	// // chat
	// http.HandleFunc("/api/Chat", chat.ChatHandler)
	// http.HandleFunc("/api/GetMessages", chat.GetMessagesHandler)
	// http.HandleFunc("/api/GetDiscussions", chat.GetDiscussionsListHandler)
	// http.HandleFunc("/api/GetOnlineUsers", chat.GetOnlineUsersHandler)

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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		fun(w, r)
	}
}
