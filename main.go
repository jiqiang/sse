package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//testtest
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		err = ws.WriteMessage(websocket.TextMessage, []byte(time.Now().Format(time.ANSIC)))
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("ui/")))
	http.HandleFunc("/ws", serveWs)

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
