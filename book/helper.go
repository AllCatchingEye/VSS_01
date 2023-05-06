package book

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

type informationHelper struct {
	books        []Book
	requestsOpen map[*actor.PID]bool
	bookActors   map[uint32]*actor.PID
	sender       *actor.PID
}

func (state *informationHelper) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case GetInformationHelper:
		state.sender = ctx.Sender()
		sender := ctx.Sender()
		println(sender)
		state.requestsOpen = make(map[*actor.PID]bool)
		for id, bookPid := range state.bookActors {
			state.requestsOpen[bookPid] = true
			ctx.Request(state.bookActors[id], GetInformationOfBook{})
		}
		fmt.Println("Helper: Requested all information")
	case Information:
		sender := msg.actorPID
		delete(state.requestsOpen, sender)
		state.books = append(state.books, msg.response)

		fmt.Println("Helper: Received information")
		if len(state.requestsOpen) == 0 {
			fmt.Println("Helper: Received all requested information")
			ctx.Send(state.sender, state.books)
		}
	}
}

// #####################################
// #        Messages for Helper        #
// #####################################

// GetInformationHelper message to collect information about all books for helper actor
type GetInformationHelper struct {
}
