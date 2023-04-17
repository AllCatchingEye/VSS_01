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
	case ServiceInformationRequest:
		helper := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &informationActor{bookActors: state.bookActors, service: ctx.Self(), client: msg.client}
		}))
		ctx.Send(helper, HelperInformationRequest{})
	case HelperInformationCollected:
		ctx.Send(msg.client, ServiceInformationCollected{books: msg.books})
	case Borrow:
		val, ok := state.bookActors[msg.id]
		if ok {
			ctx.Send(val, msg)
		} else {
			//TODO: Error Message or Poison?
			ctx.Respond(UnknownBook{})
		}
	case Return:
		val, ok := state.bookActors[msg.id]
		if ok {
			ctx.Send(val, msg)
		} else {
			//TODO: Error Message or Poison?
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
