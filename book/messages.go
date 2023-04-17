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
	books []book
}

type HelperInformationRequest struct {
}

type HelperInformationCollected struct {
	books  []book
	client *actor.PID
}

type ActorInformationResponse struct {
	response book
}

type ActorInformationRequest struct{}
