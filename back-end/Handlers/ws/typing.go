package ws

import "fmt"

func Typinsend(msg Tyoping) {
	for _, conn := range Connections[msg.Receiver] {
		err := conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}


