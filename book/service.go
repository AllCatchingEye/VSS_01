package book

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

// represents the book service with all its books
type bookServiceActor struct {
	bookActors map[uint32]*actor.PID
}

func (state *bookServiceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		newBook := Book{
			id:        1,
			title:     "Worm",
			author:    []string{"Wildbow"},
			available: 2,
			borrowed:  3,
		}
		bookActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &bookActor{book: newBook}
		}))
		state.bookActors[newBook.id] = bookActor
		fmt.Println("Book Service: Initialized")
	case GetInformation:
		// returns information of all registered books
		helper := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &informationHelper{bookActors: state.bookActors}
		}))
		ctx.RequestWithCustomSender(helper, GetInformationHelper{}, ctx.Sender())
		fmt.Println("Book Service: Dispatch information helper")
	case ReturnBook:
		// checks if book exists and requests 'Borrow'
		bookId, bookExists := state.bookActors[msg.Id]
		if bookExists {
			fmt.Println("Book Service: Borrow book from actor")
			ctx.RequestWithCustomSender(bookId, ReturnBook{Id: msg.Id, ClientId: msg.ClientId}, ctx.Sender())
		} else {
			fmt.Println("Book Service: Book not known")
			ctx.Respond(UnknownBook{})
		}
	case BorrowBook:
		// checks if book exists and requests 'Return'
		bookId, bookExists := state.bookActors[msg.Id]
		if bookExists {
			ctx.RequestWithCustomSender(bookId, BorrowBook{Id: msg.Id, ClientId: msg.ClientId}, ctx.Sender())
			fmt.Println("Book Service: Return book to actor")
		} else {
			fmt.Println("Book Service: Book not known")
			ctx.Respond(UnknownBook{})
		}
	case NewBook:
		_, bookExists := state.bookActors[msg.Book.id]
		if !bookExists {
			newActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &bookActor{book: msg.Book}
			}))
			state.bookActors[msg.Book.id] = newActor
			fmt.Println("Book Service: Added new book")
			ctx.Respond(true)
		} else {
			fmt.Println("Book Service: Coudnt add new book, book already exists with given id")
			ctx.Respond(false)
		}

	default:
		fmt.Println("got unknown message of type %T\n", msg)
	}
}

func NewBookService() actor.Actor {
	return &bookServiceActor{}
}

// #####################################
// #     Messages for Book Service     #
// #####################################

// GetInformation message to collect information about all books
type GetInformation struct {
}

// Information message holding information about single book
type Information struct {
	response Book
}

// UnknownBook message if book does not exist
type UnknownBook struct {
}

// NotAvailable message that wanted book is not available (all copies borrowed)
type NotAvailable struct {
}

// NewBook message to add a new book what will spawn new BookActor
type NewBook struct {
	Book Book
}
