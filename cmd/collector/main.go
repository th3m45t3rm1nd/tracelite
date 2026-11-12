package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/th3m45t3rm1nd/tracelite/internal/collector"
)

func main() {
	addr := os.Getenv("TRACELITE_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	c := collector.New()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/traces", c.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("tracelite collector listening on %s (buffered spans: %d)", addr, c.Len())
	log.Fatal(srv.ListenAndServe())
}
