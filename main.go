package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)
var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer ws.Close()
		clients[ws] = true
		for {
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				log.Printf("error: %v", err)
				delete(clients, ws)
				break
			}
			for client := range clients {
				client.WriteMessage(messageType, p)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
