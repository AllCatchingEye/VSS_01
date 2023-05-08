package main

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
	"time"
)

func main() {
	system := actor.NewActorSystem()

	rootContext := system.Root
	timeout := 5 * time.Second

	bsProp := actor.PropsFromProducer(book.NewBookService)
	bs := system.Root.Spawn(bsProp)

	// Test create book successfully
	bookNoExemplars := messages.Book{
		Id:        42,
		Author:    []string{"Nobody"},
		Title:     "Void",
		Available: 0,
		Borrowed:  0,
	}
	res, err := rootContext.RequestFuture(bs, &messages.NewBook{Book: &bookNoExemplars}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok := res.(*messages.BookCreated)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.BookCreated{}))
	}

	fmt.Println("Successfully created book")

	// Test Borrow book
	res, err = rootContext.RequestFuture(bs, &messages.Borrow{ClientId: 1, BookId: 1}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.(*messages.Book)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Book{}))
	}

	// Test Borrow too much
	res, err = rootContext.RequestFuture(bs, &messages.Borrow{ClientId: 1, BookId: 42}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.(*messages.NotAvailable)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.NotAvailable{}))
	}

	// Test Borrow non-existent book
	res, err = rootContext.RequestFuture(bs, &messages.Borrow{ClientId: 1, BookId: 5134}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.(*messages.UnknownBook)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.UnknownBook{}))
	}

	// Test return successfully
	res, err = rootContext.RequestFuture(bs, &messages.Return{ClientId: 1, BookId: 1}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.(*messages.Returned)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Returned{}))
	}

	// Test return too much
	res, err = rootContext.RequestFuture(bs, &messages.Return{ClientId: 1, BookId: 42}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.(*messages.NotAvailable)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.NotAvailable{}))
	}

	// Test return non-existent book
	res, err = rootContext.RequestFuture(bs, &messages.Return{ClientId: 1, BookId: 3134}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.(*messages.UnknownBook)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.UnknownBook{}))
	}

	// Test get information
	res, err = rootContext.RequestFuture(bs, &messages.GetInformation{}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.([]*messages.Book)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be []*messages.Book"))
	}

	println("All book tests successfull")
}
