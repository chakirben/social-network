package groups

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type IdOFgroup struct {
	Groupid int `json:"groupid"`
}

func GetPostGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}
	fmt.Println(userID)
	var Idgroup IdOFgroup
	err = json.NewDecoder(r.Body).Decode(&Idgroup)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}
	if Idgroup.Groupid < 1 {
		http.Error(w, "Invalid request :(", http.StatusBadRequest)
		return
	}

	
}
