package main

import (
	"fmt"
)

type event struct {
	name  string
	value float64
}

type subscription struct {
	userReg   *userRegistration
	eventName string
}

type subscriptionHub struct {
	userSubscriptions map[*userRegistration]map[string]bool
	mailingList       map[string]map[*userRegistration]bool
	subscribe         chan *subscription
	unsubscribe       chan *subscription
	events            chan *event
	stop              chan bool
}

func newSubscriptionHub() *subscriptionHub {
	return &subscriptionHub{
		userSubscriptions: make(map[*userRegistration]map[string]bool),
		mailingList:       make(map[string]map[*userRegistration]bool),
		subscribe:         make(chan *subscription, 10),
		unsubscribe:       make(chan *subscription, 10),
		events:            make(chan *event, 10),
		stop:              make(chan bool)}
}

func (subHub *subscriptionHub) run() {
	for {
		select {
		case msg := <-subHub.subscribe:

			//Logic for keeping track of user subscriptions
			subscriptions, ok := subHub.userSubscriptions[msg.userReg]
			if !ok {
				subscriptions = make(map[string]bool)
				subHub.userSubscriptions[msg.userReg] = subscriptions
			}
			subscriptions[msg.eventName] = true

			//Logic for keeping track of who is subscribed to particular
			//event
			subHub.mailingList[msg.eventName][msg.userReg] = true

			//Add logic to just run a guy to populate database at some interval
			//if an event for it doesnt exist yet
			//go populateNonsense(msg.name)

		case msg := <-subHub.unsubscribe:

			//Logic for keeping track to user subscribtions
			subscriptions, ok := subHub.userSubscriptions[msg.userReg]
			if ok {
				delete(subscriptions, msg.eventName)
			}

			//Logic for keeping track  of who is subscribed to particular
			//event
			usersSubscribed, ok := subHub.mailingList[msg.eventName]
			if ok {
				delete(usersSubscribed, msg.userReg)
			}

		case c := <-subHub.events:
			//In here means we received a message from the mongo tailer
			fmt.Println("omg I got an event")
			fmt.Println(c)
			subscribers, ok := subHub.mailingList[c.name]
			if ok {
				for subscriber, registered := range subscribers {
					if registered {
						subscriber.userWS.WriteJSON(c)
					}
				}
			}

		case <-subHub.stop:
			//Maybe kick off all users on this worker?
			return
		}
	}
}
