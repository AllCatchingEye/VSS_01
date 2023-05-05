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
		return &library.LibraryService{BookService: bs, CustomerService: cs}
	})
	ls := system.Root.Spawn(lsProp)

	// Test add customer successfully
	name := "Alice"

	res, err := rootContext.RequestFuture(ls, library.LibAddCustomer{Name: name}, timeout).Result()

	if err != nil {
		panic(err)
	}

	resCustomer, ok := res.(messages.Customer)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Customer{}))
	}

	if resCustomer.GetName() != name {
		panic(fmt.Errorf("new customer with ID %d has name %s, should be %s",
			resCustomer.GetId(), resCustomer.GetName(), name))
	}

	// Test add customer error
	name = "Alice"

	res, err = rootContext.RequestFuture(ls, library.LibAddCustomer{Name: name}, timeout).Result()
	res, err = rootContext.RequestFuture(ls, library.LibAddCustomer{Name: name}, timeout).Result()
	res, err = rootContext.RequestFuture(ls, library.LibAddCustomer{Name: name}, timeout).Result()

	if err != nil {
		panic(err)
	}

	resCustomer, ok = res.(messages.Customer)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be %T", messages.Customer{}))
	}

	if res != false {
		panic(fmt.Errorf("New customer shoudn't be successfully created"))
	}

	// Test NewBook successfully
	sut := book.CreateNewBook(2, []string{"Alice", "Hubert"}, "Super Ratgeber", 10, 0)

	resBook, err := rootContext.RequestFuture(ls, book.NewBook{Book: sut}, timeout).Result()

	if err != nil {
		panic(err)
	}

	res, ok = resBook.(bool)

	if !ok {
		panic(fmt.Errorf("go wrong message type. Should be %T", book.Book{}))
	}

	if res != true {
		panic(fmt.Errorf("coudn't create new book"))
	}

	// Test NewBook error
	resBook, err = rootContext.RequestFuture(ls, book.NewBook{Book: sut}, timeout).Result()

	if err == nil {
		panic(err)
	}

	res, ok = resBook.(bool)

	if !ok {
		panic(fmt.Errorf("go wrong message type. Should be %T", book.Book{}))
	}

	if res != false {
		panic(fmt.Errorf("could create book that already exists"))
	}

	// Test BorrowBook successfully
	var bookID uint32
	bookID = 1
	res, err = rootContext.RequestFuture(ls, book.BorrowBook{
		ClientId: 1,
		Id:       bookID,
	}, timeout).Result()

	if err != nil {
		panic(err)
	}

	borrowedBook, ok := res.(book.Book)

	if !ok {
		panic(fmt.Errorf("go wrong message type. Should be %T", book.Book{}))
	}

	if borrowedBook.GetId() != 1 {
		panic(fmt.Errorf("new book with ID %d, should be %d",
			borrowedBook.GetId(), sut.GetId()))
	}

	// Test BorrowBook error
	var borrowedBookID uint32
	borrowedBookID = 12
	res, err = rootContext.RequestFuture(ls, book.BorrowBook{
		ClientId: 1,
		Id:       borrowedBookID,
	}, timeout).Result()

	if err != nil {
		panic(err)
	}

	borrowedErrorResponse, ok := res.(bool)

	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be bool"))
	}

	if borrowedErrorResponse != false {
		panic(fmt.Errorf("book with id %d could be borrowed successfully but shouldn't", bookID))
	}

	// Test ReturnBook successfully
	res, err = rootContext.RequestFuture(ls, book.ReturnBook{
		ClientId: bookID,
		Id:       1,
	}, timeout).Result()

	if err != nil {
		panic(err)
	}

	returnResponse, ok := res.(bool)
	if !ok {
		panic(fmt.Errorf("got wrong message type. Should be bool"))
	}

	if returnResponse != true {
		panic(fmt.Errorf("book with id %d couldn't be returned successfully", bookID))
	}

	// Test ReturnBook error
	res, err = rootContext.RequestFuture(ls, book.ReturnBook{
		ClientId: bookID,
		Id:       12,
	}, timeout).Result()

	if err != nil {
		panic(err)
	}

	returnErrorResponse, ok := res.(bool)

	if !ok {
		panic(fmt.Errorf("go wrong message type. Should be %T", book.Book{}))
	}

	if returnErrorResponse != false {
		panic(fmt.Errorf("book with id %d could be returned successfully but shouldn't", bookID))
	}

}
