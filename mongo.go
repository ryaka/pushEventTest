package main

import (
	"fmt"
	"github.com/rwynn/gtm"
	"gopkg.in/mgo.v2"
	// "reflect"
)

func filterLog(op *gtm.Op) bool {
	return op.IsInsert()
}

func listenToOpLog(subHubs []*subscriptionHub) {
	//Address of where the database for mongo is located
	//in this case it is on the same machine
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	ops, errs := gtm.Tail(session, &gtm.Options{
		Filter: filterLog,
	})

	for {
		select {
		case err := <-errs:
			fmt.Println(err)
		case op := <-ops:
			name, nameOk := op.Data["name"].(string)
			value, valueOk := op.Data["value"].(float64)

			if nameOk && valueOk {
				for _, subHub := range subHubs {
					subHub.events <- &event{
						Name:  name,
						Value: value,
					}
				}
			} else {
				fmt.Println("name: ", nameOk, " value: ", valueOk)
				fmt.Println(op.Data)
			}
		}
	}
}
