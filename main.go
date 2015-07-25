package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
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
	subHubs            = make([]*subscriptionHub, 0)
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//Add this conn to a list that a wsMessageHandler will
	//handle for message subscriptions
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("shit error upgrading to websocket")
		return
	}
	fmt.Println("okay we got a websocket connection")

	userSubHub := subHubs[rand.Intn(len(subHubs))]

	userReg := &userRegistration{
		userWS: conn,
		subHub: userSubHub,
		stop:   make(chan bool),
	}

	go userReg.run()

	userSubHub.subscribe <- &subscription{
		userReg: userReg,
	}
}

func main() {
	//Spawn a certain number of "workers" that deal with subscriptions
	for i := 0; i < 1; i++ {
		subHub := newSubscriptionHub()
		go subHub.run()
		subHubs = append(subHubs, subHub)
	}
	go listenToOpLog(subHubs)
	http.HandleFunc("/ws", wsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
