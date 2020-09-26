package main

import (
	"log"
	"net/http"

	"gitlab.com/gustavocd/hellowasm/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/execute", handler.CompileCodeHandler)

	log.Println("Server is running on port :8080")
	http.ListenAndServe(":8080", mux)
}
