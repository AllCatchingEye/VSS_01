package main

import (
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {
	// run this programm forever
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	// start remote actor system
	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 9011)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	// spawn the customer service
	props := actor.PropsFromProducer(book.NewBookService)
	_, _ = system.Root.SpawnNamed(props, "bookServiceActor")
}
