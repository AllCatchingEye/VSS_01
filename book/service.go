package book

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

type bookServiceActor struct {
	bookActors map[uint32]*actor.PID
}

func (state *bookServiceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		newBook := Book{
			id:        1,
			title:     "Worm",
			author:    []string{"Wildbow"},
			available: 2,
			borrowed:  3,
		}
		bookActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &bookActor{book: newBook}
		}))
		state.bookActors[newBook.id] = bookActor
		fmt.Println("Book Service: Initialized")
	case GetInformation:
		helper := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &informationHelper{bookActors: state.bookActors}
		}))
		ctx.RequestWithCustomSender(helper, msg, ctx.Sender())
		fmt.Println("Book Service: Dispatch information helper")
	case Borrow:
		bookId, bookExists := state.bookActors[msg.Id]
		if bookExists {
			fmt.Println("Book Service: Borrow book from actor")
			ctx.RequestWithCustomSender(bookId, msg, ctx.Sender())
		} else {
			fmt.Println("Book Service: Book not known")
			ctx.Respond(UnknownBook{})
		}
	case Return:
		bookId, bookExists := state.bookActors[msg.Id]
		if bookExists {
			ctx.RequestWithCustomSender(bookId, msg, ctx.Sender())
			fmt.Println("Book Service: Return book to actor")
		} else {
			fmt.Println("Book Service: Book not known")
			ctx.Respond(UnknownBook{})
		}
	case NewBook:
		newActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &bookActor{book: msg.Book}
		}))
		state.bookActors[msg.Book.id] = newActor
		fmt.Println("Book Service: Added new book")
	default:
		fmt.Println("got unknown message of type %T\n", msg)
	}
}
