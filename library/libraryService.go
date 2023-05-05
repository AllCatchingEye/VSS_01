package library

import (
	"fmt"

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
		fmt.Println("Library Service: New customer to add")
	case book.NewBook:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransNewBook{
			bookService:     state.bookService,
			customerService: state.customerService,
			book:            msg.Book,
		}, ctx.Sender())
		fmt.Println("Library Service: New book to add")
	case book.Borrow:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransBorrow{
			bookService:     state.bookService,
			customerService: state.customerService,
			bookMsg:         msg,
		}, ctx.Sender())
		fmt.Println("Library Service: Book borrow requested")
	case book.Return:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransReturn{
			bookService:     state.bookService,
			customerService: state.customerService,
			bookMsg:         msg,
		}, ctx.Sender())
		fmt.Println("Library Service: Book return requested")
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
