package library

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
	"time"
)

type transActor struct {
}

func (state *transActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case TransAddCustomer:
		authClient(ctx, msg.customerService, ctx.Sender().)
		f := ctx.RequestFuture(msg.customerService, messages.NewCustomer{Name: msg.name}, 50*time.Millisecond)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			ctx.Send(ctx.Parent(), res)
		} else {
			ctx.Respond(err)
			ctx.Send(ctx.Parent(), err)
		}
	case TransNewBook:
		f := ctx.RequestFuture(msg.bookService, book.NewBook{Book: msg.book}, 50*time.Millisecond)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			ctx.Send(ctx.Parent(), res)
		} else {
			ctx.Respond(err)
			ctx.Send(ctx.Parent(), err)
		}
	case TransBorrow:
		f := ctx.RequestFuture(msg.bookService, msg.bookMsg, 50*time.Millisecond)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			ctx.Send(ctx.Parent(), res)
		} else {
			ctx.Respond(err)
			ctx.Send(ctx.Parent(), err)
		}
	case TransReturn:
		f := ctx.RequestFuture(msg.bookService, msg.bookMsg, 50*time.Millisecond)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			ctx.Send(ctx.Parent(), res)
		} else {
			ctx.Respond(err)
			ctx.Send(ctx.Parent(), err)
		}
	default:
		print("Error occured by handling following message: %T\n", msg)
	}
}

func authClient(ctx actor.Context, customerService *actor.PID, clientId uint32) bool {
	f := ctx.RequestFuture(customerService, messages.GetCustomer{Id: clientId}, 50*time.Millisecond)
	res, err := f.Result()
	if err != nil {
		return false
	}
	switch res.(type) {
	case messages.CustomerNotFound:
		return false
	default:
		return true
	}
}
