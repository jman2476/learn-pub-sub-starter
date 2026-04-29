package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jman2476/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jman2476/learn-pub-sub-starter/internal/pubsub"
	"github.com/jman2476/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()

	connString := "amqp://guest:guest@localhost:5672"
	connection, err := amqp.Dial(connString)
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error establishing connection to RabbitMQ: %w", err),
		)
	}
	defer connection.Close()
	fmt.Println("Connection to RabbitMQ server successful!")

	channel, err := connection.Channel()
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error creating channel: %w", err),
		)
	}

	_, logQueue, err := pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		strings.Join([]string{routing.GameLogSlug, "*"}, "."),
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error creating game_logs queue: %w", err),
		)
	}
	log.Printf("Queue: %v", logQueue)

	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}

		switch inputs[0] {
		case "pause":
			log.Println("Sending pause message")
			err = pubsub.PublishJSON(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
			if err != nil {
				log.Println(
					fmt.Errorf("Error pausing game: %w", err),
				)
			}
		case "resume":
			log.Println("Sending resume message")
			err = pubsub.PublishJSON(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
			if err != nil {
				log.Println(
					fmt.Errorf("Error resuming game: %w", err),
				)
			}
		case "quit":
			log.Println("Exiting server, goodbye")
			return
		default:
			log.Println("I didn't understand that command")
		}

	}
}
