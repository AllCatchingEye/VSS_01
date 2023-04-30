package book

import "github.com/asynkron/protoactor-go/actor"

type informationHelper struct {
	books        []Book
	requestsOpen uint32
	bookActors   map[uint32]*actor.PID
}

func (state informationHelper) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case GetInformation:
		state.requestsOpen = 0
		for id, _ := range state.bookActors {
			state.requestsOpen += 1
			ctx.Request(state.bookActors[id], msg)
		}
	case Information:
		state.requestsOpen -= 1
		state.books = append(state.books, msg.response)

		if state.requestsOpen < 1 {
			ctx.Respond(state.books)
			ctx.Poison(ctx.Self())
		}
	}
}
