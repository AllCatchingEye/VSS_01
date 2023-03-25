package main

import (
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"gitlab.lrz.de/vss/startercode/startercodeB1/customer"
)

func main() {
	// run this programm forever
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	// start remote actor system
	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 9010)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	// spawn the customer service
	props := actor.PropsFromProducer(customer.NewService)
	_, _ = system.Root.SpawnNamed(props, "customer")
}
