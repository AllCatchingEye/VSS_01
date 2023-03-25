package customer

import (
	"github.com/asynkron/protoactor-go/actor"

	"gitlab.lrz.de/vss/startercode/startercodeB1/messages"
)

type service struct {
	customer map[uint32]*customer
	lastID   uint32
}

func NewService() actor.Actor {
	return &service{customer: make(map[uint32]*customer), lastID: 0}
}

func (cs *service) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.NewCustomer:
		name := msg.GetName()
		id := cs.addNewCustomer(name)

		ctx.Respond(&messages.Customer{Name: name, Id: id})

	case *messages.GetCustomer:
		id := msg.GetId()
		customer, ok := cs.getCustomer(id)

		if ok {
			ctx.Respond(&messages.Customer{Name: customer.name, Id: customer.id})
		} else {
			ctx.Respond(&messages.CustomerNotFound{Id: id})
		}
	}
}

func (cs *service) addNewCustomer(name string) uint32 {
	cs.lastID++
	id := cs.lastID
	cs.customer[id] = &customer{id: id, name: name}

	return id
}

func (cs *service) getCustomer(id uint32) (*customer, bool) {
	c, ok := cs.customer[id]

	return c, ok
}
