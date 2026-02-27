package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Server struct {
	Addr    string
	Handler *http.ServeMux
}

type apiConfig struct {
	fileserverHits atomic.Int32
	Handler        *http.ServeMux
}

func main() {
	var apiCfg apiConfig
	mux := http.NewServeMux()
	server := &Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.Handler.HandleFunc("/healthz/", server.ServeHTTP)
	server.Handler.HandleFunc("/metrics/", apiCfg.middlewareHits)
	server.Handler.Handle(
		"/app/",
		http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))),
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

func (api *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (api *apiConfig) middlewareHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Hits: %d\n", api.fileserverHits.Load())
	if err != nil {
		fmt.Println("failed to write response counter")
	}
}
