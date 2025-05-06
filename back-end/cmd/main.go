package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"socialN/cmd/routers"

	db "socialN/dataBase"
)

func main() {
	db.DbInit()
	routers.SetupHandlers()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Println("your serve on : http://localhost:8080")

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("server error:", err)
		}
	}()

	<-sig
	
	err := db.SocialDB.Close()
	if err != nil {
		log.Println("\nerror to close database", err)
	} else {
		fmt.Println("\ndatabase is closed (:")
	}
}
