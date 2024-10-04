package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/olahol/melody"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		// read message from client
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			break
		}
		// show message
		log.Printf("Received message: %s", message)

		//send message to client
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			break
		}

	}
}

func main() {
	useNative := false
	if useNative {
		http.HandleFunc("/websocket", websocketHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		m := melody.New()
		m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		fs := http.FileServer(http.Dir("./"))
		http.Handle("/", fs)
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			m.HandleRequest(w, r)
		})
		m.HandleMessage(func(s *melody.Session, msg []byte) {
			m.Broadcast(msg)
		})
		log.Println("Server started at :3000")
		log.Fatal(http.ListenAndServe(":3000", nil))
	}

}
