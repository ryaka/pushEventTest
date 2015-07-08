package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type wsMessage struct {
	conn        *websocket.Conn
	messageType int
	p           []byte
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	eventSubscriptions = map[string][]*websocket.Conn{}
	subHub             = newSubscriptionHub()
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//Add this conn to a list that a wsMessageHandler will
	//handle for message subscriptions
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			//log.Println(err)
			return
		}
		fmt.Println(messageType, p)
	}
}

func main() {
	//TODO launch a goroutine that listens in on oplog and then
	//sends to the list of subscribers
	go subHub.run()
	go listenToOpLog(subHub)
	http.HandleFunc("/ws", wsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
