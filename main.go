package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	Addr    string
	Handler http.Handler
}

func main() {
	server := Server{Addr: ":8080"}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(server.Addr, mux)
	if err != nil {
		fmt.Println("failed to start server")
	}
}
