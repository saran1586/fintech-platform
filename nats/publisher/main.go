package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	err = nc.Publish(
	"payments.initiated.PAY001",
	[]byte("hello"),
)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message published")
}