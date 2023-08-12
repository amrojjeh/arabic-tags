package main

import (
	"net/http"
	"os"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

func main() {
	mux := http.NewServeMux()
	s := http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	defer s.Close()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
