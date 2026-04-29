package main

import (
	"fmt"
	"log"
	"os"
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

	channel, err := connection.Channel()
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error creating channel: %w", err),
		)
	}

	var username string
	if len(os.Args) > 1 && os.Args[1] != "" {
		username, err = gamelogic.QuickClientWelcome(os.Args[1])
	} else {
		username, err = gamelogic.ClientWelcome()
	}
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error getting username: %w", err),
		)
	}

	gameState := gamelogic.NewGameState(username)

	// Subscribe to PAUSE
	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilDirect,
		strings.Join([]string{routing.PauseKey, username}, "."),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gameState),
	)
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error subscribing to queue: %w", err),
		)
	}

	// Subscribe to ARMY_MOVES.*
	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilTopic,
		strings.Join([]string{routing.ArmyMovesPrefix, username}, "."),
		strings.Join([]string{routing.ArmyMovesPrefix, "*"}, "."),
		pubsub.SimpleQueueTransient,
		handlerMove(gameState),
	)

	// REPL start
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
			move, err := gameState.CommandMove(inputs)
			if err != nil {
				fmt.Println(
					fmt.Errorf("Error moving units: %w", err),
				)
				continue
			}

			err = pubsub.PublishJSON(
				channel,
				routing.ExchangePerilTopic,
				strings.Join([]string{routing.ArmyMovesPrefix, username}, "."),
				move,
			)
			if err != nil {
				fmt.Println(
					fmt.Errorf("Error publishing move: %w", err),
				)
				continue
			}
			log.Printf("Move %v successfully published", move)

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
