package main

import ()

type event struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
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
	unsubscribeAll    chan *subscription
	events            chan *event
	stop              chan bool
}

func newSubscriptionHub() *subscriptionHub {
	return &subscriptionHub{
		userSubscriptions: make(map[*userRegistration]map[string]bool),
		mailingList:       make(map[string]map[*userRegistration]bool),
		subscribe:         make(chan *subscription, 10),
		unsubscribe:       make(chan *subscription, 10),
		unsubscribeAll:    make(chan *subscription, 10),
		events:            make(chan *event, 10),
		stop:              make(chan bool)}
}

func (subHub *subscriptionHub) unsubscribeUserAll(msg *subscription) {
	//Logic for keeping track to user subscribtions
	subscriptions, ok := subHub.userSubscriptions[msg.userReg]
	if !ok {
		return
	}

	for event := range subscriptions {
		subHub.unsubscribeUserEvent(&subscription{
			userReg:   msg.userReg,
			eventName: event,
		})
	}

	delete(subHub.userSubscriptions, msg.userReg)
}

func (subHub *subscriptionHub) unsubscribeUserEvent(msg *subscription) {
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
			if msg.eventName != "" {
				subscriptions[msg.eventName] = true

				//Logic for keeping track of who is subscribed to particular
				//event
				usersSubscribed, ok := subHub.mailingList[msg.eventName]
				if !ok {
					usersSubscribed = make(map[*userRegistration]bool)
					subHub.mailingList[msg.eventName] = usersSubscribed
				}
				subHub.mailingList[msg.eventName][msg.userReg] = true

				//Add logic to just run a guy to populate database at some interval
				//if an event for it doesnt exist yet
				//go populateNonsense(msg.name)
			}

		case msg := <-subHub.unsubscribe:
			subHub.unsubscribeUserEvent(msg)
		case msg := <-subHub.unsubscribeAll:
			subHub.unsubscribeUserAll(msg)
		case c := <-subHub.events:
			//In here means we received a message from the mongo tailer
			subscribers, ok := subHub.mailingList[c.Name]
			if ok {
				for subscriber, registered := range subscribers {
					if registered {
						err := subscriber.userWS.WriteJSON(c)

						if err != nil {
							subHub.unsubscribeUserAll(&subscription{
								userReg: subscriber,
							})
						}

					}
				}
			}

		case <-subHub.stop:
			//Maybe kick off all users on this worker?
			return
		}
	}
}
