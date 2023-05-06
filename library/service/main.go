package main

import (
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/library"
)

func main() {
	// run this programm forever
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	// start remote actor system
	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 9012)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	// spawn the customer service
	props := actor.PropsFromProducer(library.NewLibraryService)
	_, _ = system.Root.SpawnNamed(props, "LibraryService")
}
