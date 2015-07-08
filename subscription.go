package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type event struct {
	name  string
	value float64
}

type subscription struct {
	user *websocket.Conn
	name string
	all  bool
}

type subscriptionHub struct {
	userSubscriptions map[*websocket.Conn][]string
	mailingList       map[string][]*websocket.Conn
	subscribe         chan *subscription
	unsubscribe       chan *subscription
	events            chan *event
	stop              chan bool
}

func newSubscriptionHub() *subscriptionHub {
	return &subscriptionHub{
		userSubscriptions: make(map[*websocket.Conn][]string),
		mailingList:       make(map[string][]*websocket.Conn),
		subscribe:         make(chan *subscription, 10),
		unsubscribe:       make(chan *subscription, 10),
		events:            make(chan *event, 10),
		stop:              make(chan bool)}
}

func (subHub *subscriptionHub) run() {
	for {
		select {
		case <-subHub.subscribe:
		case <-subHub.unsubscribe:
		case c := <-subHub.events:
			fmt.Println("omg I got an event")
			fmt.Println(c)
		case <-subHub.stop:
			return
		}
	}
}
