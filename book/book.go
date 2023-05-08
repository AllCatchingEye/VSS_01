package book

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
)

// represents a book
type bookActor struct {
	book *messages.Book
}

func (state *bookActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.GetInformation:
		fmt.Println("Book Actor: Information requested")
		ctx.Respond(InformationFound{book: state.book, sender: ctx.Self()})
	case *messages.Borrow:
		if state.book.Available > 0 {
			state.book.Available -= 1
			state.book.Borrowed += 1
			fmt.Println("Book Actor: Book borrowed")
			ctx.Respond(state.book)
		} else {
			fmt.Println("Book Actor: Coudn't borrow, no book available right now")
			ctx.Respond(&messages.NotAvailable{})
		}
	case *messages.Return:
		if state.book.Borrowed > 0 {
			state.book.Available += 1
			state.book.Borrowed -= 1
			fmt.Println("Book Actor: Book returned")
			ctx.Respond(&messages.Returned{})
		} else {
			fmt.Println("Book Actor: Coudn't return, no book was borrowed")
			ctx.Respond(&messages.NotAvailable{})
		}
	default:
		fmt.Printf("got a message of type %T\n", msg)
	}
}

type InformationFound struct {
	book   *messages.Book
	sender *actor.PID
}
