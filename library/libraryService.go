package library

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
)

type LibraryService struct {
	BookService           *actor.PID
	CustomerService       *actor.PID
	transActors           []*actor.PID
	transActorMap         map[*actor.PID]bool // true: transactor is running
	transActorCustomerMap map[*actor.PID]uint32
	runningTransActors    int
}

func (state *LibraryService) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.transActorMap = make(map[*actor.PID]bool)
		state.transActorCustomerMap = make(map[*actor.PID]uint32)
	case *messages.LibAddCustomer:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx, 0), TransAddCustomer{
			name:            msg.Name,
			customerService: state.CustomerService,
		}, ctx.Sender())
		fmt.Println("Library Service: New customer to add")
	case *messages.NewBook:
		ctx.RequestWithCustomSender(state.spawnTransActor(ctx, 0), TransNewBook{
			bookService:     state.BookService,
			customerService: state.CustomerService,
			newBookMessage:  msg,
		}, ctx.Sender())
		fmt.Println("Library Service: New newBookMessage to add")
	case *messages.Borrow:
		customer := msg.ClientId
		if state.checkTransactorClientMap(customer) {
			ctx.RequestWithCustomSender(state.spawnTransActor(ctx, customer), TransBorrow{
				bookService:     state.BookService,
				customerService: state.CustomerService,
				borrowMessage:   msg,
			}, ctx.Sender())
			fmt.Println("Library Service: Book borrow requested")
		} else {
			ctx.Respond(&messages.SameCustomer{})
		}
	case *messages.Return:
		customer := msg.ClientId
		if state.checkTransactorClientMap(customer) {
			ctx.RequestWithCustomSender(state.spawnTransActor(ctx, customer), TransReturn{
				bookService:     state.BookService,
				customerService: state.CustomerService,
				returnMessage:   msg,
			}, ctx.Sender())
			fmt.Println("Library Service: Book return requested")
		} else {
			ctx.Respond(&messages.SameCustomer{})
		}
	case TransFinished:
		fmt.Println("Entering library received transfinished")
		state.runningTransActors--
		state.transActorMap[ctx.Sender()] = false // not running anymore
		delete(state.transActorCustomerMap, msg.sender)
		fmt.Println("deleted customer from map")
	default:
		print("Unknown message. %T\n", msg)
	}
}

// spawns new Transactor or returns existing and free transactor
func (state *LibraryService) spawnTransActor(ctx actor.Context, customer uint32) *actor.PID {
	// clearng exessive transactors
	inactiveTransactors := len(state.transActors) - state.runningTransActors
	if inactiveTransactors >= 3 {
		state.clearExcessiveTransactors()
	}
	// find transactor in existing transactors that is free and return its PID
	for _, PIDid := range state.transActors {
		if !state.transActorMap[PIDid] {
			state.transActorMap[PIDid] = true // running again
			if customer != 0 {
				state.transActorCustomerMap[PIDid] = customer
			}
			state.runningTransActors++
			return PIDid
		}
	}
	// If no free transactor is found create a new one and append it to transActors
	trans := ctx.Spawn(actor.PropsFromProducer(NewTransActor))
	state.transActors = append(state.transActors, trans)
	state.transActorMap[trans] = true // running
	if customer != 0 {
		state.transActorCustomerMap[trans] = customer
	}
	state.runningTransActors++
	return trans
}

// deletes all execcive transactors
func (state *LibraryService) clearExcessiveTransactors() {
	for _, PIDid := range state.transActors {
		if state.transActorMap[PIDid] {
			delete(state.transActorMap, PIDid)
			delete(state.transActorCustomerMap, PIDid)
		}
	}
}

// checks if there is a transactor working for this client
func (state *LibraryService) checkTransactorClientMap(customerId uint32) bool {
	for _, trans := range state.transActors {
		if state.transActorCustomerMap[trans] == customerId {
			return false
		}
	}
	return true
}

func NewLibraryService(bs *actor.PID, cs *actor.PID) actor.Actor {
	println("Creating new library service")
	return &LibraryService{BookService: bs, CustomerService: cs}
}

type TransFinished struct {
	sender *actor.PID
}
