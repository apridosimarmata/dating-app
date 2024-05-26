package main

import (
	"dating-app/presentation"
	"fmt"
	"net/http"
)

func main() {
	router := presentation.InitServer()

	http.ListenAndServe(":3000", router)
	fmt.Println("server listening on port 3000")
}
