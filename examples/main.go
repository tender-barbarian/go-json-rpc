package main

import (
	"go-json-rpc/examples/api"
	"go-json-rpc/pkg/rpc"
	"log"
)

func main() {
	server, err := rpc.New(":8080")
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
