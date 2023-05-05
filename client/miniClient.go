package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/customer"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/library"
	"time"
)

func main() {
	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 0)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	cs := actor.NewPID("127.0.0.1:9010", "customer")
	bs := actor.NewPID("127.0.0.1:9011", "bookServiceActor")

	libraryProps := actor.PropsFromProducer(func() actor.Actor {
		return NewLibraryService(bs, cs)
	})
	ls := actor.Spawn(libraryProps)

	system := actor.NewActorSystem()

	rootContext := system.Root
	timeout := 5 * time.Second

	csProp := actor.PropsFromProducer(customer.NewService)
	cs := system.Root.Spawn(csProp)

	bsProp := actor.PropsFromProducer(book.NewBookService)
	bs := system.Root.Spawn(bsProp)

	lsProp := actor.PropsFromProducer(func() actor.Actor {
		return &library.LibraryService{BookService: bs, CustomerService: cs}
	})
	ls := system.Root.Spawn(lsProp)

	res, err := rootContext.RequestFuture(ls, book.BorrowBook{
		ClientId: 1,
		Id:       1,
	}, timeout).Result()

}
