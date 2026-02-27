package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	Addr    string
	Handler *http.ServeMux
}

func main() {
	mux := http.NewServeMux()
	server := &Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.Handler.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(server.Addr, server.Handler)
	if err != nil {
		fmt.Println("failed to start server")
	}
}
