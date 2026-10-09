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

	_, err = nc.Subscribe("test", func(msg *nats.Msg) {
		fmt.Println("Received:", string(msg.Data))
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Subscriber is listening...")

	select {}
}