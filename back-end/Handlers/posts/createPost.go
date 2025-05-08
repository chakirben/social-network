package posts

import (
	"fmt"
	"net/http"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in it ")
}
