package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/bye", handler2)

	// srv := http.Server{
	// 	Addr:    ":3000",
	// 	Handler: mux,
	// }

	// fmt.Println("Starting auth server")
	// if err := srv.ListenAndServe(); err != nil {
	// 	fmt.Printf("server auth: %s\n", err)
	// }
	_ = http.ListenAndServe(":8081", mux)
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("auth"))
}

func handler2(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("auth bye"))
}
