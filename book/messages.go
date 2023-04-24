package book

import "github.com/asynkron/protoactor-go/actor"

type Borrow struct {
	Client *actor.PID
	Id     uint32
}

type Return struct {
	Client *actor.PID
	Id     uint32
}

type ServiceInformationRequest struct {
	client *actor.PID
}

type ServiceInformationCollected struct {
	books []Book
}

type HelperInformationRequest struct {
}

type HelperInformationCollected struct {
	books  []Book
	client *actor.PID
}

type ActorInformationResponse struct {
	response Book
}

type ErrorBook struct {
}

type ActorInformationRequest struct{}

type NewBook struct {
	Book Book
}
