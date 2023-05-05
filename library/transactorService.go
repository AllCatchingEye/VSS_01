package library

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
)

type transActor struct {
}

func (state *transActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Trans actor: Initilized")
	case TransAddCustomer:
		f := ctx.RequestFuture(msg.customerService, messages.NewCustomer{Name: msg.name}, 50*time.Millisecond)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			ctx.Send(ctx.Parent(), res)
			fmt.Println("Trans Actor: Added customer")
		} else {
			ctx.Respond(err)
			ctx.Send(ctx.Parent(), err)
			fmt.Println("Trans Actor: Error: Something went wrong while trying to add customer")
		}
	case TransNewBook:
		f := ctx.RequestFuture(msg.bookService, book.NewBook{Book: msg.book}, 50*time.Millisecond)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			ctx.Send(ctx.Parent(), res)
			fmt.Println("Trans Actor: Added new book")
		} else {
			ctx.Respond(err)
			ctx.Send(ctx.Parent(), err)
			fmt.Println("Trans Actor: Error: Something went wrong while trying to add book")
		}
	case TransBorrow:
		authClient(ctx, msg.customerService, msg.bookMsg.ClientId)
		f := ctx.RequestFuture(msg.bookService, msg.bookMsg, 50*time.Millisecond)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			ctx.Send(ctx.Parent(), res)
			fmt.Println("Trans Actor: Borrowed book")
		} else {
			ctx.Respond(err)
			ctx.Send(ctx.Parent(), err)
			fmt.Println("Trans Actor: Error: Something went wrong while trying to borrow book")
		}
	case TransReturn:
		authClient(ctx, msg.customerService, msg.bookMsg.ClientId)
		f := ctx.RequestFuture(msg.bookService, msg.bookMsg, 50*time.Millisecond)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			ctx.Send(ctx.Parent(), res)
			fmt.Println("Trans Actor: Returned book")
		} else {
			ctx.Respond(err)
			ctx.Send(ctx.Parent(), err)
			fmt.Println("Trans Actor: Error: Something went wrong while trying to return book")
		}
	default:
		print("Error occured by handling following message: %T\n", msg)
	}
}

func authClient(ctx actor.Context, customerService *actor.PID, clientId uint32) bool {
	f := ctx.RequestFuture(customerService, messages.GetCustomer{Id: clientId}, 50*time.Millisecond)
	res, err := f.Result()
	if err != nil {
		fmt.Println("Trans Actor: Something went wrong while trying to check for customer")
		return false
	}
	switch res.(type) {
	case messages.CustomerNotFound:
		fmt.Println("Trans Actor: Customer not found")
		return false
	default:
		return true
	}
}
