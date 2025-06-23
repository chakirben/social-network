package app

import (
	"fmt"
	"net/http"
	"strconv"

	"socialN/app/auth"

	database "socialN/dataBase"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, database.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	itemType := r.URL.Query().Get("itemType")
	itemId := r.URL.Query().Get("itemId")
	reaction_type, err := strconv.Atoi(r.URL.Query().Get("reactionType"))
	if err != nil {
		http.Error(w, "Invalid reaction type, it should be 1 or 0, -1", http.StatusBadRequest)
		return

	}

	if itemType == "" || itemId == "" {
		http.Error(w, "Missing item Type or item id", http.StatusBadRequest)
		return
	}
	fmt.Println(itemType, itemId)
	if itemType != "post" && itemType != "comment" {
		http.Error(w, "Invalid item type", http.StatusBadRequest)
		return
	}
	fmt.Println(itemId, itemType, reaction_type)
	if itemType == "post" {
		_, err := database.SocialDB.Exec(`INSERT INTO postReactions (postId, userId, reactionType) VALUES (?, ?, ?)
		ON CONFLICT(postId, userId) 
		DO UPDATE SET reactionType = excluded.reactionType`, itemId, userID, reaction_type)
		if err != nil {
			http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Reaction added successfully"))
	}
	if itemType == "comment" {
		_, err := database.SocialDB.Exec(`INSERT INTO commentReactions (commentId, userId, reactionType) VALUES (?, ?, ?)
		ON CONFLICT(commentId, userId) 
		DO UPDATE SET reactionType = excluded.reactionType`, itemId, userID, reaction_type)
		if err != nil {
			http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Reaction added successfully"))
		return
	}
}
