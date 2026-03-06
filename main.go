package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filePathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz/", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp/", handlerResponse)
	mux.HandleFunc("GET /admin/metrics/", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset/", apiCfg.handlerReset)
	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))),
	)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())

}
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
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

func (api *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, `<html>
					<body>
						<h1>Welcome, Chirpy Admin</h1>
						<p>Chirpy has been visited %d times!</p>
					</body>
				</html>`, api.fileserverHits.Load())
	if err != nil {
		fmt.Println("failed to write response counter")
	}
}

func (api *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	api.fileserverHits.Store(0)
	_, err := fmt.Fprintf(w, "Hits: %d\n", api.fileserverHits.Load())
	if err != nil {
		fmt.Println("failed to reset counter")
	}
}

func handlerResponse(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Body string `json:"body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	var body requestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Something went wrong"})
		return
	}

	if len(body.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Chirp is too long"})
		return
	}

	words := replaceBadWords(body.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(validResponse{CleanedBody: words})
}

func replaceBadWords(s string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(s, " ")
	for i, word := range words {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}
