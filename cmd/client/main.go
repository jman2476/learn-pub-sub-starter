package main

import (
	"fmt"
	"strings"

	"github.com/jman2476/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jman2476/learn-pub-sub-starter/internal/pubsub"
	"github.com/jman2476/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	connstring := "amqp://guest:guest@localhost:5672"
	connection, err := amqp.Dial(connstring)
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error establishing RabbitMQ connection: %w", err),
		)
	}
	defer connection.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error getting username: %w", err),
		)
	}

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilDirect,
		strings.Join([]string{routing.PauseKey, username}, "."),
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gameState),
	)
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error subscribing to queue: %w", err),
		)
	}

	for {
		inputs := gamelogic.GetInput()

		if len(inputs) == 0 {
			continue
		}

		switch inputs[0] {
		case "spawn":
			err := gameState.CommandSpawn(inputs)
			if err != nil {
				fmt.Println(
					fmt.Errorf("Error spawing units: %w", err),
				)
			}
		case "move":
			_, err := gameState.CommandMove(inputs)
			if err != nil {
				fmt.Println(
					fmt.Errorf("Error moving units: %w", err),
				)
				continue
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command. Use the 'help' to see possible commands")
		}
	}

}
