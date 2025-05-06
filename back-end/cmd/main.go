package main

import (
	"fmt"
	"net/http"

	"socialN/cmd/routers"

	db "socialN/dataBase"
)

func main() {
	db.DbInit()
	routers.SetupHandlers()
	fmt.Println("https://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
