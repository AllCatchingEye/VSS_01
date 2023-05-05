package book

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

type Book struct {
	id        uint32
	author    []string
	title     string
	available uint32
	borrowed  uint32
}

// represents a book
type bookActor struct {
	book Book
}

func (state *bookActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case GetInformationOfBook:
		fmt.Println("Book Actor: Information requested")
		ctx.Respond(Information{response: state.book})
	case BorrowBook:
		if state.book.available > 0 {
			state.book.available -= 1
			state.book.borrowed += 1
			fmt.Println("Book Actor: Book borrowed")
			ctx.Respond(state.book)
		} else {
			fmt.Println("Book Actor: Coudn't borrow, no book available right now")
			ctx.Respond(NotAvailable{})
		}
	case ReturnBook:
		if state.book.borrowed > 0 {
			state.book.available += 1
			state.book.borrowed -= 1
			fmt.Println("Book Actor: Book returned")
			ctx.Respond(true)
		} else {
			fmt.Println("Book Actor: Coudn't return, no book was borrowed")
			ctx.Respond(false)
		}
	default:
		fmt.Printf("got a message of type %T\n", msg)
	}
}

func CreateNewBook(id uint32, author []string, title string, available uint32, borrowed uint32) Book {
	return Book{
		id:        id,
		author:    author,
		title:     title,
		available: available,
		borrowed:  borrowed,
	}
}

func (x *Book) GetId() uint32 {
	if x != nil {
		return x.id
	}
	return 0
}

func (x *Book) GetTitle() string {
	if x != nil {
		return x.title
	}
	return ""
}

func (x *Book) GetAuthors() []string {
	if x != nil {
		return x.author
	}
	return []string{}
}

// #####################################
// #         Messages for Book         #
// #####################################

// BorrowBook message to borrow book
type BorrowBook struct {
	ClientId uint32
	Id       uint32
}

// ReturnBook message to return book
type ReturnBook struct {
	ClientId uint32
	Id       uint32
}

// GetInformationOfBook message to collect information about this book
type GetInformationOfBook struct {
}
