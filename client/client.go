package main

import (
	"fmt"
	emitter "github.com/emitter-io/go/v2"
)

type RelayClient struct {
	EmitterClient *emitter.Client
}

func main() {
	fmt.Printf("Client")
}
