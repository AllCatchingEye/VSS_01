package book

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
)

type book struct {
	id        uint32
	author    []string
	title     string
	available uint32
	borrowed  uint32
}

type bookActor struct {
	book book
}

func (state *bookActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case ActorInformationRequest:
		ctx.Respond(ActorInformationResponse{response: state.book})
	case Borrow:
		if state.book.borrowed > 0 {
			state.book.borrowed -= 1
			ctx.Send(msg.client, state.book)
		} else {
			ctx.Send(msg.client, &actor.PoisonPill{})
		}
	case Return:
		state.book.available += 1
		ctx.Send(msg.client, state.book)
	default:
		fmt.Printf("got a message of type %T\n", msg)
	}
}
