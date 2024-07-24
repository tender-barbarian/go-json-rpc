package main

import (
	"go-json-rpc/internal/api"
	"go-json-rpc/server"
	"log"
)

func main() {
	server, err := server.New(":8080")
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
