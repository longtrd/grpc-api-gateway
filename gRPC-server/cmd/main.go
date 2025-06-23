package main

import (
	"log"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/cmd/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}
