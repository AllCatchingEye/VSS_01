package book

import "github.com/asynkron/protoactor-go/actor"

type informationActor struct {
	books        []Book
	requestsOpen uint32
	service      *actor.PID
	client       *actor.PID
	bookActors   map[uint32]*actor.PID
}

func (state informationActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case HelperInformationRequest:
		state.requestsOpen = 0
		for id, _ := range state.bookActors {
			state.requestsOpen += 1
			ctx.Request(state.bookActors[id], ActorInformationRequest{})
		}
		//TODO: Warten bis alle den Request bekommen haben?
	case ActorInformationResponse:
		state.requestsOpen -= 1
		state.books = append(state.books, msg.response)

		if state.requestsOpen < 1 {
			ctx.Send(state.service, HelperInformationCollected{books: state.books, client: state.client})
			ctx.Poison(ctx.Self())
		}
	}
}
