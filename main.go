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
	server.Handler.HandleFunc("/healthz/", server.ServeHTTP)
	server.Handler.Handle(
		"/app/",
		http.StripPrefix("/app", http.FileServer(http.Dir("."))),
	)
	err := http.ListenAndServe(server.Addr, server.Handler)
	if err != nil {
		fmt.Println("failed to start server")
	}
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println("failed to write response")
	}
}
