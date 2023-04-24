package library

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
)

type TransAddCustomer struct {
	name            string
	customerService *actor.PID
}

type TransNewBook struct {
	bookService     *actor.PID
	customerService *actor.PID
	book            book.Book
}

type TransBorrow struct {
	bookService     *actor.PID
	customerService *actor.PID
	bookMsg         book.Borrow
}

type TransReturn struct {
	bookService     *actor.PID
	customerService *actor.PID
	bookMsg         book.Return
}
