package main

import (
	"fmt"
	"io"
	"net/http"
)

func WatchdogHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	header := w.Header()
	header.Set("Content-Type", "text/plain")
	io.WriteString(w, "OK")
}

func MetricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	header := w.Header()
	header.Set("Content-Type", "text/html")
	x := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %v times!</p></body></html>", Config.MetricsGet())
	w.Write([]byte(x))
}

func MetricsResetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	Config.MetricsReset()
	err := Config.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		fmt.Printf("Error deleting users from db:%s", err)
		w.Write([]byte("Error deleting users from db"))
	}
}

func (cfg *apiconfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiconfig) MetricsGet() int32 {
	return cfg.fileserverHits.Load()
}

func (cfg *apiconfig) MetricsReset() {
	cfg.fileserverHits.Swap(0)
}
