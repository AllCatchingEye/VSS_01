package main

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
)

func main() {
	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 0)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	service := actor.NewPID("127.0.0.1:9010", "customer")

	name := "Max Mustermann"
	res, err := system.Root.RequestFuture(
		service,
		&messages.NewCustomer{Name: name},
		5*time.Second).Result()

	if err != nil {
		panic(err)
	}

	customer, ok := res.(*messages.Customer)

	if !ok {
		panic(fmt.Errorf("go wrong message type"))
	}

	fmt.Printf(">>> got answer: %v\n", customer)
}
