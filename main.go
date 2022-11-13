package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	startServer()
}

func startServer() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request)
		go handleHandShake()
	})
	log.Fatal(http.ListenAndServe(":6000", nil))
}

func handleHandShake() {
      
}
func listenForMessages()  {
	
}