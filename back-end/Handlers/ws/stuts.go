package ws

import (
	"fmt"
)

func sendStutus(msg Status) {
	fmt.Println("Sending status message:", msg)
	for _, users := range Connections {
		for _, conn1 := range users {
			err := conn1.WriteJSON(msg)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
			fmt.Println("rah mxa")
		}
	}
}
