package chiserver

import (
	"fmt"
	"net/http"
)

func StartServer() {
	fmt.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", routes())
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
