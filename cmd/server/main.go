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

	serverChannel, _, err := pubsub.DeclareAndBind(dial, routing.ExchangePerilDirect, "game_logs", routing.GameLogSlug, "durable")
	if err != nil {
		return
	}
	defer serverChannel.Close()

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if input[0] == "pause" {
			fmt.Println("Sending pause message to clients...")
			err = pubsub.PublishJSON(newChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				return
			}
		} else if input[0] == "resume" {
			fmt.Println("Sending resume message to clients...")
			err = pubsub.PublishJSON(newChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				return
			}
		} else if input[0] == "quit" {
			fmt.Println("Exiting...")
			break
		} else {
			fmt.Println("Unknown command.")
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nExiting...")
}
