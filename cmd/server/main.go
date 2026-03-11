package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionString := "amqp://guest:guest@localhost:5672/"
	dial, err := amqp.Dial(connectionString)
	if err != nil {
		return
	}
	defer dial.Close()

	fmt.Println("Connection Successful!")

	newChannel, err := dial.Channel()
	if err != nil {
		return
	}
	defer newChannel.Close()

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nExiting...")

	pubsub.PublishJSON(newChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
}
