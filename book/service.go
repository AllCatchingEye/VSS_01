package book

import (
	"fmt"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"

	"github.com/asynkron/protoactor-go/actor"
)

// represents the book service with all its books
type bookServiceActor struct {
	bookActors map[uint32]*actor.PID
}

func (state *bookServiceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		newBook := &messages.Book{
			Id:        1,
			Author:    []string{"Wildbow"},
			Title:     "Worm",
			Available: 2,
			Borrowed:  3,
		}
		bookActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &bookActor{book: newBook}
		}))
		state.bookActors = make(map[uint32]*actor.PID)
		state.bookActors[newBook.Id] = bookActor
		fmt.Println("Book Service: Initialized")
	case *messages.GetInformation:
		// returns information of all registered books
		helper := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &informationHelper{bookActors: state.bookActors}
		}))
		ctx.RequestWithCustomSender(helper, &messages.GetInformation{}, ctx.Sender())
		fmt.Println("Book Service: Dispatch information helper")
	case *messages.Return:
		// checks if book exists and requests 'Borrow'
		bookId, bookExists := state.bookActors[msg.BookId]
		if bookExists {
			fmt.Println("Book Service: Borrow book from actor")
			ctx.RequestWithCustomSender(bookId, &messages.Return{BookId: msg.BookId, ClientId: msg.ClientId}, ctx.Sender())
		} else {
			fmt.Println("Book Service: Book not known")
			ctx.Respond(&messages.UnknownBook{})
		}
	case *messages.Borrow:
		// checks if book exists and requests 'Return'
		bookId, bookExists := state.bookActors[msg.BookId]
		if bookExists {
			ctx.RequestWithCustomSender(bookId, &messages.Borrow{BookId: msg.BookId, ClientId: msg.ClientId}, ctx.Sender())
			fmt.Println("Book Service: Return book to actor")
		} else {
			fmt.Println("Book Service: Book not known")
			ctx.Respond(&messages.UnknownBook{})
		}
	case *messages.NewBook:
		_, bookExists := state.bookActors[msg.Book.Id]
		if !bookExists {
			newActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &bookActor{book: msg.Book}
			}))
			state.bookActors[msg.Book.Id] = newActor
			fmt.Println("Book Service: Added new book")
			ctx.Respond(&messages.BookCreated{})
		} else {
			fmt.Println("Book Service: Coudnt add new book, book already exists with given id")
			ctx.Respond(&messages.BookExists{})
		}
	default:
		print("Unknown message. %T\n", msg)
	}
}

func NewBookService() actor.Actor {
	return &bookServiceActor{}
}
