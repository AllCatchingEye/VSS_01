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

	// Test Borrow book
	res, err := rootContext.RequestFuture(bs, book.BorrowBook{ClientId: 1, Id: 1}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok := res.(book.Book)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Customer{}))
	}

	// Test Borrow too much
	rootContext.RequestFuture(bs, book.BorrowBook{ClientId: 1, Id: 1}, timeout).Result()
	res, err = rootContext.RequestFuture(bs, book.BorrowBook{ClientId: 1, Id: 1}, timeout).Result()

	if err != nil {
		panic(err)
	}

	resBook, ok := res.(bool)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Customer{}))
	}

	if resBook {
		panic(fmt.Errorf("Shoudn't be able to borrow when the book is unavailable"))
	}

	// Test Borrow non-existent book
	res, err = rootContext.RequestFuture(bs, book.BorrowBook{ClientId: 1, Id: 5134}, timeout).Result()

	if err != nil {
		panic(err)
	}

	resBook, ok = res.(bool)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Customer{}))
	}

	if resBook {
		panic(fmt.Errorf("Shoudn't be able to borrow book"))
	}

	// Test return successfully
	res, err = rootContext.RequestFuture(bs, book.ReturnBook{ClientId: 1, Id: 1}, timeout).Result()

	if err != nil {
		panic(err)
	}

	resBook, ok = res.(bool)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Customer{}))
	}

	if !resBook {
		panic(fmt.Errorf("Should be able to return book"))
	}

	// Test return too much
	rootContext.RequestFuture(bs, book.ReturnBook{ClientId: 1, Id: 1}, timeout).Result()
	rootContext.RequestFuture(bs, book.ReturnBook{ClientId: 1, Id: 1}, timeout).Result()
	rootContext.RequestFuture(bs, book.ReturnBook{ClientId: 1, Id: 1}, timeout).Result()
	rootContext.RequestFuture(bs, book.ReturnBook{ClientId: 1, Id: 1}, timeout).Result()
	rootContext.RequestFuture(bs, book.ReturnBook{ClientId: 1, Id: 1}, timeout).Result()
	res, err = rootContext.RequestFuture(bs, book.ReturnBook{ClientId: 1, Id: 1}, timeout).Result()

	if err != nil {
		panic(err)
	}

	resBook, ok = res.(bool)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Customer{}))
	}

	if resBook {
		panic(fmt.Errorf("Should be able to return book that has all exemplars already available"))
	}

	// Test return non-existent book
	res, err = rootContext.RequestFuture(bs, book.ReturnBook{ClientId: 1, Id: 3134}, timeout).Result()

	if err != nil {
		panic(err)
	}

	resBook, ok = res.(bool)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Customer{}))
	}

	if resBook {
		panic(fmt.Errorf("Should be able to return book that doesnt exist"))
	}

	// Test get information
	res, err = rootContext.RequestFuture(bs, book.GetInformation{}, timeout).Result()

	if err != nil {
		panic(err)
	}

	resInfo, ok := res.([]book.Book)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", book.Book{}))
	}

	if len(resInfo) != 1 {
		panic(fmt.Errorf("Should be able to return book that doesnt exist"))
	}
	if resInfo[0].GetId() != 1 {
		panic(fmt.Errorf("book coming from info should have id %d, but has id %d", 1, resInfo[0].GetId()))
	}
	if resInfo[0].GetTitle() != "Worm" {
		panic(fmt.Errorf("book coming from info should have title %s, but has title %s", "Worm", resInfo[0].GetTitle()))
	}

	println("All book tests successfull")
}
