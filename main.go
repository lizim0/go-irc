package main

import (
	"container/list"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{}

func main() {
	l := list.New()
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
		for e := l.Front(); e != nil; e = e.Next() {
			ws.WriteMessage(websocket.TextMessage, e.Value.([]byte))
		}
		for {
			_, p, err := ws.ReadMessage()
			if err != nil {
				log.Printf("error: %v", err)
				delete(clients, ws)
				break
			}
			l.PushBack(p)
			for client := range clients {
				client.WriteMessage(websocket.TextMessage, p)
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
