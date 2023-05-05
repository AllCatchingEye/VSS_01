package library

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/book"
)

type LibraryService struct {
	BookService        *actor.PID
	CustomerService    *actor.PID
	transActors        []*actor.PID
	transActorMap      map[*actor.PID]bool // true: transactor is running
	runningTransActors int
}

func (state *LibraryService) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case LibAddCustomer:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransAddCustomer{
			name:            msg.Name,
			customerService: state.CustomerService,
		}, ctx.Sender())
		fmt.Println("Library Service: New customer to add")
	case book.NewBook:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransNewBook{
			bookService:     state.BookService,
			customerService: state.CustomerService,
			book:            msg.Book,
		}, ctx.Sender())
		fmt.Println("Library Service: New book to add")
	case book.BorrowBook:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransBorrow{
			bookService:     state.BookService,
			customerService: state.CustomerService,
			bookMsg:         msg,
		}, ctx.Sender())
		fmt.Println("Library Service: Book borrow requested")
	case book.ReturnBook:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx), TransReturn{
			bookService:     state.BookService,
			customerService: state.CustomerService,
			bookMsg:         msg,
		}, ctx.Sender())
		fmt.Println("Library Service: Book return requested")
	case bool:
		if msg {
			state.runningTransActors--
			state.transActorMap[ctx.Sender()] = false // not running anymore
		} else {
			fmt.Printf("Transactor %s: task failure\n", ctx.Sender().String())
		}
	default:
		print("Unknown message. %T\n", msg)
	}
}

// spawns new Transactor or returns existing and free transactor
func (state *LibraryService) spawnTransActor(ctx actor.Context) *actor.PID {
	// clearng exessive transactors
	inactiveTransactors := len(state.transActors) - state.runningTransActors
	if inactiveTransactors >= 3 {
		state.clearExcessiveTransactors()
	}
	// find transactor in existing transactors that is free and return its PID
	for _, PIDid := range state.transActors {
		if !state.transActorMap[PIDid] {
			state.transActorMap[PIDid] = true // running again
			state.runningTransActors++
			return PIDid
		}
	}
	// If no free transactor is found create a new one and append it to transActors
	trans := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return &transActor{}
	}))
	state.transActors = append(state.transActors, trans)
	state.transActorMap[trans] = true // running
	state.runningTransActors++
	return trans
}

// deletes all execcive transactors
func (state *LibraryService) clearExcessiveTransactors() {
	for _, PIDid := range state.transActors {
		if state.transActorMap[PIDid] {
			delete(state.transActorMap, PIDid)
		}
	}
}

// #####################################
// #       Messages for Library        #
// #####################################

// TransAddCustomer message for transactor to add a new customer
type LibAddCustomer struct {
	Name string
}
