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

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(dial, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, "transient", gameState.HandlePause)
	if err != nil {
		return
	}

	for {
		input := gamelogic.GetInput()
		locations := []string{"americas", "europe", "asia", "africa", "australia", "antarctica"}
		units := []string{"infantry", "cavalry", "artillery"}
		if input[0] == "spawn" {
			for _, loc := range locations {
				if loc == input[1] {
					for _, unit := range units {
						if unit == input[2] {
							err := gameState.CommandSpawn(input)
							if err != nil {
								return
							}
						}
					}
				}
			}
		} else if input[0] == "move" {
			_, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println("Invalid move.")
			} else {
				fmt.Println("Move successful.")
			}
		} else if input[0] == "status" {
			gameState.CommandStatus()
		} else if input[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if input[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		} else if input[0] == "quit" {
			gamelogic.PrintQuit()
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
