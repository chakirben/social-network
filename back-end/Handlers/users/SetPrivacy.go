package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

func SetPrivacy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqBody struct {
		AccountType string `json:"accountType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if reqBody.AccountType != "private" && reqBody.AccountType != "public" {
		http.Error(w, "Invalid accountType", http.StatusBadRequest)
		return
	}

	_, err = dataB.SocialDB.Exec("UPDATE Users SET accountType = ? WHERE id = ?", reqBody.AccountType, userID)
	if err != nil {
		fmt.Println("error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Privacy updated successfully"}`))
}
