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
				fmt.Println("I sent", name, value)
				for _, subHub := range subHubs {
					subHub.events <- &event{
						name:  name,
						value: value,
					}
				}
			} else {
				fmt.Println("name: ", nameOk, " value: ", valueOk)
				fmt.Println(op.Data)
			}

			// msg := fmt.Sprintf(`Got op <%v> for object <%v>
			// in database <%v>
			// and collection <%v>
			// and data <%v>
			// and timestamp <%v>`,
			// 	op.Operation, op.Id, op.GetDatabase(),
			// 	op.GetCollection(), op.Data, op.Timestamp)
			// fmt.Println(msg) // or do something more interesting
			// fmt.Println(fmt.Sprintf(`Operation: %v`, reflect.TypeOf(op.Operation)))
			// fmt.Println(fmt.Sprintf(`Id: %v`, reflect.TypeOf(op.Id)))
			// fmt.Println(fmt.Sprintf(`GetDatabase: %v`, reflect.TypeOf(op.GetDatabase())))
			// fmt.Println(fmt.Sprintf(`GetCollection: %v`, reflect.TypeOf(op.GetCollection())))
			// fmt.Println(fmt.Sprintf(`Data: %v`, reflect.TypeOf(op.Data)))
			// fmt.Println(fmt.Sprintf(`Timestamp: %v`, reflect.TypeOf(op.Timestamp)))
		}
	}
}
