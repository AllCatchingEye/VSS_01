package book

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

type informationHelper struct {
	books        []Book
	requestsOpen uint32
	bookActors   map[uint32]*actor.PID
}

func (state informationHelper) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.requestsOpen = 0
	case GetInformationHelper:
		for id := range state.bookActors {
			state.requestsOpen += 1
			ctx.Request(state.bookActors[id], GetInformationOfBook{})
		}
		fmt.Println("Helper: Requested all information")
	case InformationForHelper:
		state.requestsOpen -= 1
		state.books = append(state.books, msg.response)

		fmt.Println("Helper: Received information")
		if state.requestsOpen < 1 {
			fmt.Println("Helper: Received all requested information")
			ctx.Respond(state.books)
			ctx.Poison(ctx.Self())
		}
	}
}

// #####################################
// #        Messages for Helper        #
// #####################################

// GetInformationHelper message to collect information about all books for helper actor
type GetInformationHelper struct {
}

// InformationForHelper message holding information about single book for helper actor
type InformationForHelper struct {
	response Book
}
