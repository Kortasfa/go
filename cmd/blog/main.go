package main

import (
	"log"
	"net/http"
)

const (
	port = ":3000"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/home", index)
	mux.HandleFunc("/post", post)

	// Реализуем отдачу статики
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Start server " + port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
