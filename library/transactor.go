package library

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
	"time"
)

type transActor struct {
	// TODO: dürfen wir hier einen boolean speichern um zu wissen ob der actor fertig ist oder noch arbeitet
}

func (state *transActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	// TODO: wie antworten wir hier dem libraryService und dem Costumer/Client
	case TransBorrow:
		f := ctx.RequestFuture(msg.bookService, msg.bookMsg, 50*time.Millisecond)
		res, err := f.Result()
		if err != nil {
			ctx.Respond(res)
		} else {
			ctx.Respond(err)
		}
	case TransReturn:
		f := ctx.RequestFuture(msg.bookService, msg.bookMsg, 50*time.Millisecond)
		res, err := f.Result()
		if err != nil {
			ctx.Respond(res)
		} else {
			ctx.Respond(err)
		}
	case TransAddBook:
		f := ctx.RequestFuture(msg.bookService, msg, 50*time.Millisecond)
		res, err := f.Result()
		if err != nil {
			ctx.Respond(res)
		} else {
			ctx.Respond(err)
		}
	case TransAddCustomer:
		f := ctx.RequestFuture(msg.customerService, messages.NewCustomer{Name: msg.name}, 50*time.Millisecond)
		res, err := f.Result()
		if err != nil {
			ctx.Respond(res)
		} else {
			ctx.Respond(err)
		}
	}
}
