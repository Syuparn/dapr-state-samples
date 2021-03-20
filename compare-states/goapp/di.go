package main

import (
	dapr "github.com/dapr/go-sdk/client"
	"go.uber.org/dig"
)

func NewContainer() (*dig.Container, func()) {
	teardowns := []func(){}

	c := dig.New()

	// client
	c.Provide(func() dapr.Client {
		cli, err := dapr.NewClientWithPort("50001") // default
		if err != nil {
			panic(err)
		}
		teardowns = append(teardowns, func() { cli.Close() })
		return cli
	})

	// repository
	c.Provide(NewMessageRepository)

	// inputPort
	c.Provide(NewCreateMessagePort)

	// controller
	c.Provide(NewController)

	teardown := func() {
		for _, t := range teardowns {
			t()
		}
	}

	return c, teardown
}
