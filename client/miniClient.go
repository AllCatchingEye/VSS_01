package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
	"time"
)

func main() {
	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 9013)
	remoter := remote.NewRemote(system, config)
	remoter.Start()
	rootContext := system.Root
	timeout := 5 * time.Second

	ls := actor.NewPID("127.0.0.1:9012", "LibraryService")

	res, err := rootContext.RequestFuture(ls, &messages.Borrow{ClientId: 1, BookId: 1}, timeout).Result()

	if err != nil {
		panic("Something went wrong while trying to borrow a book")
	}

	book, ok := res.(*messages.Book)
	if !ok {
		println("Coudn't borrow book")
	}

	println("Got book: ", book.Title)
}
