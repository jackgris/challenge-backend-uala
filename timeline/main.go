package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	srv := http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	fmt.Println("Starting timeline server")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("server timeline: %s\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("timeline"))
}
