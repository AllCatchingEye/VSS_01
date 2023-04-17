package library

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
)

type libraryActor struct {
	bookService     *actor.PID
	customerService *actor.PID
	transActors     []*actor.PID
	transActorMap   map[*actor.PID]bool
}

func (state *libraryActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	// TODO: Wie prüfen wir ob ein customer bereits existiert?
	case book.Borrow:
		ctx.Send(state.spawnTransActor(ctx), msg)
	case book.Return:
		ctx.Send(state.spawnTransActor(ctx), msg)
	case book.NewBook:
		ctx.Send(state.spawnTransActor(ctx), msg)
	case *messages.NewCustomer:
		ctx.Send(state.spawnTransActor(ctx), msg)
	}
}

func (state *libraryActor) spawnTransActor(ctx actor.Context) *actor.PID {
	transActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return &transActor{}
	}))
	state.transActors = append(state.transActors, transActor)
	return transActor
}
