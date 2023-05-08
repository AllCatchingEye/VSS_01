package library

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
)

type transActor struct {
}

func (state *transActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Trans actor: Initialized")
	case TransAddCustomer:
		f := ctx.RequestFuture(msg.customerService, &messages.NewCustomer{Name: msg.name}, 5*time.Second)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			fmt.Println("Trans Actor: Added customer")
		} else {
			ctx.Respond(err)
			fmt.Println("Trans Actor: Error: Something went wrong while trying to add customer")
		}
		ctx.Send(ctx.Parent(), TransFinished{ctx.Self()})
	case TransNewBook:
		f := ctx.RequestFuture(msg.bookService, msg.newBookMessage, 5*time.Second)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			fmt.Println("Trans Actor: Added new newBookMessage")
		} else {
			ctx.Respond(err)
			fmt.Println("Trans Actor: Error: Something went wrong while trying to add newBookMessage")
		}
		ctx.Send(ctx.Parent(), TransFinished{ctx.Self()})
	case TransBorrow:
		fmt.Println("Entering trans borrow")
		authClient(ctx, msg.customerService, msg.borrowMessage.ClientId)
		f := ctx.RequestFuture(msg.bookService, msg.borrowMessage, 5*time.Second)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			fmt.Println("Trans Actor: Borrowed newBookMessage")
		} else {
			ctx.Respond(err)
			fmt.Println("Trans Actor: Error: Something went wrong while trying to borrow newBookMessage")
		}
		ctx.Send(ctx.Parent(), TransFinished{ctx.Self()})
		fmt.Println("Finished trans borrow")
	case TransReturn:
		authClient(ctx, msg.customerService, msg.returnMessage.ClientId)
		f := ctx.RequestFuture(msg.bookService, msg.returnMessage, 5*time.Second)
		res, err := f.Result()
		if err == nil {
			ctx.Respond(res)
			fmt.Println("Trans Actor: Returned newBookMessage")
		} else {
			ctx.Respond(err)
			fmt.Println("Trans Actor: Error: Something went wrong while trying to return newBookMessage")
		}
		ctx.Send(ctx.Parent(), TransFinished{ctx.Self()})
	default:
		print("Error occured by handling following message: %T\n", msg)
	}
}

func authClient(ctx actor.Context, customerService *actor.PID, clientId uint32) bool {
	f := ctx.RequestFuture(customerService, &messages.GetCustomer{Id: clientId}, 5*time.Second)
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

func NewTransActor() actor.Actor {
	return &transActor{}
}

// #####################################
// #      Messages for Transactor      #
// #####################################

// TransAddCustomer message for transactor to add a new customer
type TransAddCustomer struct {
	name            string
	customerService *actor.PID
}

// TransNewBook message for transactor to add a new newBookMessage
type TransNewBook struct {
	bookService     *actor.PID
	customerService *actor.PID
	newBookMessage  *messages.NewBook
}

// TransBorrow message for transactor to borrow a specific newBookMessage
type TransBorrow struct {
	bookService     *actor.PID
	customerService *actor.PID
	borrowMessage   *messages.Borrow
}

// TransReturn message for transactor to return a specific newBookMessage
type TransReturn struct {
	bookService     *actor.PID
	customerService *actor.PID
	returnMessage   *messages.Return
}
