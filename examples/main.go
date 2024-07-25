package main

import (
	"log"

	"github.com/tender-barbarian/go-json-rpc/examples/api"
)

func main() {
	server, err := New(":8080")
	if err != nil {
		log.Fatalf("initializing RPC server: %v", err)
		return
	}

	err = server.Handler.Register("Health.Check", api.Health{}.Check)
	if err != nil {
		log.Fatalf("registering methods: %v", err)
		return
	}

	err = server.Start()
	if err != nil {
		log.Fatalf("starting RPC server: %v", err)
		return
	}
}
