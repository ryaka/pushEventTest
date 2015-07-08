package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type userMessage struct {
	action    string
	eventName string
}

type userRegistration struct {
	userWS *websocket.Conn
	subHub *subscriptionHub
	stop   chan bool
}

func (userReg *userRegistration) run() {
	for {
		userMsg := userMessage{}
		err := userReg.userWS.ReadJSON(&userMsg)
		if err != nil {
			//log.Println(err)
			return
		}

		if userMsg.action == "subscribe" {
			userReg.subscribeTo(userMsg.eventName)
		} else if userMsg.action == "unsubscribe" {
			userReg.unsubscribeFrom(userMsg.eventName)
		}
		fmt.Println(userMsg)

		//Add logic to subscribe the user to more stuff?
		//userReg.subScribeTo(p)
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
