package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/library"
	"time"
)

func main() {
	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 0)
	remoter := remote.NewRemote(system, config)
	remoter.Start()
	rootContext := system.Root
	timeout := 5 * time.Second

	cs := actor.NewPID("127.0.0.1:9010", "customer")
	bs := actor.NewPID("127.0.0.1:9011", "bookServiceActor")
	ls := actor.NewPID("127.0.0.1:9012", "LibraryService")
	rootContext.Send(ls, library.LibAddServices{cs, bs})

	res, err := rootContext.RequestFuture(ls, book.BorrowBook{
		ClientId: 1,
		Id:       1,
	}, timeout).Result()

	if err != nil {
		panic("Something went wrong while trying to borrow a book")
	}

	borrowedBook, ok := res.(book.Book)
	if ok {
		println("Borrowed the following book: %s", borrowedBook.GetTitle())
	}
}
