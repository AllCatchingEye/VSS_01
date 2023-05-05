package book

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

type Book struct {
	id        uint32
	author    []string
	title     string
	available uint32
	borrowed  uint32
}

type bookActor struct {
	book Book
}

func (state *bookActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case GetInformation:
		fmt.Println("Book Actor: Information requested")
		ctx.Respond(Information{response: state.book})
	case Borrow:
		if state.book.available > 0 {
			state.book.available -= 1
			state.book.borrowed += 1
			fmt.Println("Book Actor: Book borrowed")
			ctx.Respond(state.book)
		} else {
			fmt.Println("Book Actor: Coudn't borrow, no book available right now")
			ctx.Respond(NotAvailable{})
		}
	case Return:
		if state.book.borrowed > 0 {
			state.book.available += 1
			state.book.borrowed -= 1
			fmt.Println("Book Actor: Book returned")
			ctx.Respond(state.book)
		} else {
			fmt.Println("Book Actor: Coudn't return, no book was borrowed")
			ctx.Respond(NotAvailable{})
		}
	default:
		fmt.Printf("got a message of type %T\n", msg)
	}
}
