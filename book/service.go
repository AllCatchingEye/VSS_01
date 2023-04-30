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
	case GetInformation:
		helper := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &informationHelper{bookActors: state.bookActors}
		}))
		ctx.RequestWithCustomSender(helper, msg, ctx.Sender())
	case Borrow:
		bookId, bookExists := state.bookActors[msg.Id]
		if bookExists {
			ctx.RequestWithCustomSender(bookId, msg, ctx.Sender())
		} else {
			ctx.Respond(UnknownBook{})
		}
	case Return:
		bookId, bookExists := state.bookActors[msg.Id]
		if bookExists {
			ctx.RequestWithCustomSender(bookId, msg, ctx.Sender())
		} else {
			ctx.Respond(UnknownBook{})
		}
	case NewBook:
		newActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &bookActor{book: msg.Book}
		}))
		state.bookActors[msg.Book.id] = newActor
	default:
		fmt.Println("got unknown message of type %T\n", msg)
	}
}
