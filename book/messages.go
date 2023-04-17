package book

import "github.com/asynkron/protoactor-go/actor"

type Borrow struct {
	client *actor.PID
	id     uint32
}

type Return struct {
	client *actor.PID
	id     uint32
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

type UnknownBook struct {
}

type ActorInformationRequest struct{}

type NewBook struct {
	Book Book
}
