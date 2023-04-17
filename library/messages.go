package library

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
)

type TransBorrow struct {
	bookMsg         book.Borrow
	bookService     *actor.PID
	customerService *actor.PID
}

type TransReturn struct {
	bookMsg         book.Return
	bookService     *actor.PID
	customerService *actor.PID
}

type TransAddBook struct {
	bookService     *actor.PID
	customerService *actor.PID
}

type TransAddCustomer struct {
	name            string
	bookService     *actor.PID
	customerService *actor.PID
}
