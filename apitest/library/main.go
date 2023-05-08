package main

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/customer"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/library"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
	"time"
)

func main() {
	system := actor.NewActorSystem()

	rootContext := system.Root
	timeout := 5 * time.Second

	csProp := actor.PropsFromProducer(customer.NewService)
	cs := system.Root.Spawn(csProp)

	bsProp := actor.PropsFromProducer(book.NewBookService)
	bs := system.Root.Spawn(bsProp)

	lsProp := actor.PropsFromProducer(func() actor.Actor {
		return library.NewLibraryService(bs, cs)
	})
	ls := system.Root.Spawn(lsProp)

	// Test add customer
	name := "Alice"

	res, err := rootContext.RequestFuture(ls, &messages.LibAddCustomer{Name: name}, timeout).Result()

	if err != nil {
		panic(err)
	}

	resCustomer, ok := res.(*messages.Customer)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Customer{}))
	}

	if resCustomer.GetName() != name {
		panic(fmt.Errorf("new customer with ID %d has name %s, should be %s",
			resCustomer.GetId(), resCustomer.GetName(), name))
	}

	// Test NewBook successfully
	sut := messages.Book{
		Id:        2,
		Author:    []string{"Alice", "Hubert"},
		Title:     "Super Ratgeber",
		Available: 10,
		Borrowed:  0,
	}

	resBook, err := rootContext.RequestFuture(ls, &messages.NewBook{Book: &sut}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = resBook.(*messages.BookCreated)

	if !ok {
		panic(fmt.Errorf("got wrong message type, expected message of type BookCreated"))
	}

	fmt.Println("Created book successfully")

	// Test Create duplicate book
	resBook, err = rootContext.RequestFuture(ls, &messages.NewBook{Book: &sut}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = resBook.(*messages.BookExists)

	if !ok {
		panic(fmt.Errorf("got wrong message type, expected message of type BookExists"))
	}

	fmt.Println("Coudn't create duplicate book!")

	// Test BorrowBook successfully
	var bookID uint32 = 1
	res, err = rootContext.RequestFuture(ls, &messages.Borrow{
		ClientId: 1,
		BookId:   bookID,
	}, timeout).Result()

	if err != nil {
		panic(err)
	}

	borrowedBook, ok := res.(*messages.Book)

	if !ok {
		panic(fmt.Errorf("got wrong message type, expected message of type Book"))
	}

	fmt.Printf("Borrowed book: %v\n", borrowedBook)

	// Test Borrow unknown book
	var borrowedBookID uint32 = 123749
	res, err = rootContext.RequestFuture(ls, &messages.Borrow{
		ClientId: 1,
		BookId:   borrowedBookID,
	}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.(*messages.UnknownBook)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be message of type UnknownBook"))
	}

	fmt.Println("Coudnt borrow unknown book!")

	// Test ReturnBook successfully
	res, err = rootContext.RequestFuture(ls, &messages.Return{
		ClientId: bookID,
		BookId:   1,
	}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.(*messages.Returned)
	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be message of type Returned"))
	}

	fmt.Println("Book returned successfully!")

	// Test return unknown book
	res, err = rootContext.RequestFuture(ls, &messages.Return{
		ClientId: bookID,
		BookId:   12,
	}, timeout).Result()

	if err != nil {
		panic(err)
	}

	_, ok = res.(*messages.UnknownBook)

	if !ok {
		panic(fmt.Errorf("got wrong message type, expected message of type UnknownBook"))
	}

	fmt.Println("Coudnt return unknown book!")

	println("All library tests successfull")
}
