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
	case *messages.NewCustomer:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransAddCustomer{
			name:            msg.Name,
			customerService: state.customerService,
		}, ctx.Sender())
	case book.NewBook:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransNewBook{
			bookService:     state.bookService,
			customerService: state.customerService,
			book:            msg.Book,
		}, ctx.Sender())
	case book.Borrow:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransBorrow{
			bookService:     state.bookService,
			customerService: state.customerService,
			bookMsg:         msg,
		}, ctx.Sender())
	case book.Return:
		ctx.Send(state.spawnTransActor(ctx), TransReturn{
			bookService:     state.bookService,
			customerService: state.customerService,
			bookMsg:         msg,
		})
	default:
		print("Unknown message. %T\n", msg)
	}
}

func (state *libraryActor) spawnTransActor(ctx actor.Context) *actor.PID {
	transActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return &transActor{}
	}))
	state.transActors = append(state.transActors, transActor)
	return transActor
}
