package book

import (
	"fmt"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type informationHelper struct {
	books        []*messages.Book
	requestsOpen map[*actor.PID]bool
	bookActors   map[uint32]*actor.PID
	sender       *actor.PID
}

func (state *informationHelper) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.GetInformation:
		state.sender = ctx.Sender()
		sender := ctx.Sender()
		println(sender)
		state.requestsOpen = make(map[*actor.PID]bool)
		for id, bookPid := range state.bookActors {
			state.requestsOpen[bookPid] = true
			ctx.Request(state.bookActors[id], msg)
		}
		fmt.Println("Helper: Requested all information")
	case InformationFound:
		delete(state.requestsOpen, msg.sender)
		state.books = append(state.books, msg.book)

		fmt.Println("Helper: Received information")
		if len(state.requestsOpen) == 0 {
			fmt.Println("Helper: Received all requested information")
			ctx.Send(state.sender, state.books)
		}
	}
}
