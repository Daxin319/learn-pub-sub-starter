package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		return
	}

	clientChannel, _, err := pubsub.DeclareAndBind(dial, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, "transient")
	if err != nil {
		return
	}
	defer clientChannel.Close()

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nExiting...")
}
