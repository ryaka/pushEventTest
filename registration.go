package main

import (
	"github.com/gorilla/websocket"
)

type userMessage struct {
	Action    string `json:"action"`
	EventName string `json:"eventName"`
}

type userRegistration struct {
	userWS *websocket.Conn
	subHub *subscriptionHub
	stop   chan bool
}

func (userReg *userRegistration) run() {
	for {
		userMsg := &userMessage{}
		err := userReg.userWS.ReadJSON(&userMsg)
		if err != nil {
			//Connection closed, unsubscribe user from everything
			userReg.unsubscribeFromAll()
			return
		}

		if userMsg.Action == "subscribe" {
			userReg.subscribeTo(userMsg.EventName)
		} else if userMsg.Action == "unsubscribe" {
			userReg.unsubscribeFrom(userMsg.EventName)
		}

	}
}

func (userReg *userRegistration) subscribeTo(eventName string) {
	userReg.subHub.subscribe <- &subscription{
		userReg:   userReg,
		eventName: eventName,
	}
}

func (userReg *userRegistration) unsubscribeFrom(eventName string) {
	userReg.subHub.unsubscribe <- &subscription{
		userReg:   userReg,
		eventName: eventName,
	}
}

func (userReg *userRegistration) unsubscribeFromAll() {
	userReg.subHub.unsubscribeAll <- &subscription{
		userReg: userReg,
	}
}
