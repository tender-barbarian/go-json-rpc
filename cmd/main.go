package main

import (
	"go-json-rpc/internal/api"
	"go-json-rpc/server"
	"log"
)

func main() {
	s, err := server.New(":8080")
	if err != nil {
		log.Fatalf("initializing RPC server: %v", err)
		return
	}

	err = s.Handler.Register("Health", &api.Health{})
	if err != nil {
		log.Fatalf("registering methods: %v", err)
		return
	}

	err = s.Start()
	if err != nil {
		log.Fatalf("starting RPC server: %v", err)
		return
	}
}
