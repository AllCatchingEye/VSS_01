//nolint:goerr113
package main

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/customer"
	"gitlab.lrz.de/vss/semester/ob-23ss/blatt-1/blatt1-grp06/messages"
)

func main() {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(customer.NewService)
	service := system.Root.Spawn(props)

	rootContext := system.Root
	timeout := 5 * time.Second

	name := "Max Mustermann"
	res, err := rootContext.RequestFuture(
		service,
		&messages.NewCustomer{Name: name},
		timeout).Result()

	if err != nil {
		panic(err)
	}

	customer, ok := res.(*messages.Customer)

	if !ok {
		panic(fmt.Errorf("go wrong message type"))
	}

	if customer.GetName() != name {
		panic(fmt.Errorf("new customer with ID %d has name %s, should be %s",
			customer.GetId(), customer.GetName(), name))
	}

	// get customer with existing id
	id := customer.GetId()
	res, err = rootContext.RequestFuture(
		service,
		&messages.GetCustomer{Id: id},
		timeout).Result()

	if err != nil {
		panic(err)
	}

	customer, ok = res.(*messages.Customer)

	if !ok {
		panic(fmt.Errorf("go wrong message type"))
	}

	if customer.GetId() != id || customer.GetName() != name {
		panic(fmt.Errorf("customer with ID %d has name %s, should be ID %d and name %s",
			customer.GetId(), customer.GetName(), id, name))
	}

	fmt.Printf(">>> got answer: %v\n", customer)

	// get ustomer with nonexisting id
	wrongID := id + 1
	res, err = rootContext.RequestFuture(
		service,
		&messages.GetCustomer{Id: wrongID},
		timeout).Result()

	if err != nil {
		panic(err)
	}

	customerNotFound, ok := res.(*messages.CustomerNotFound)

	if !ok {
		panic(fmt.Errorf("go wrong message type"))
	}

	if customerNotFound.GetId() != wrongID {
		panic(fmt.Errorf("wrong ID in customerNotFound, got %d, want %d", customerNotFound.GetId(), wrongID))
	}

	time.Sleep(time.Second)
}
