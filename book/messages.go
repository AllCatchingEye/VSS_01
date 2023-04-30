package book

import "github.com/asynkron/protoactor-go/actor"

type Borrow struct {
	Client   *actor.PID
	ClientId uint32
	Id       uint32
}

type Return struct {
	Client   *actor.PID
	ClientId uint32
	Id       uint32
}

type GetInformation struct {
	client *actor.PID
}

type Information struct {
	response Book
}

type ErrorBook struct {
}

type NotAvailable struct {
}

type NewBook struct {
	Book Book
}
