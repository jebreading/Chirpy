package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {

	var apiCfg = apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /reset", apiCfg.resetHandler)
	mux.HandleFunc("GET /metrics", apiCfg.newhitshandler)
	muxserver := http.Server{
		Addr:    ":8080",
		Handler: mux}

	err := http.ListenAndServe(muxserver.Addr, muxserver.Handler)

	if err != nil {
		fmt.Println(err) //check server started correctly.
	}

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment the counter
		cfg.fileserverHits.Add(1)
		// Call the wrapped handler
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	// Reset the counter to 0
	cfg.fileserverHits.Store(0)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Counter reset"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (hand *apiConfig) newhitshandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", hand.fileserverHits.Load())))
}
