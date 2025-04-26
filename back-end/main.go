package main

import (
	"fmt"
	"net/http"

	db "socialN/dataBase"
)

func main() {
	db.DbInit()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	fmt.Println("https://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
