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

	err = pubsub.PublishJSON(
		channel,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: true},
	)
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error publishing JSON pause: %w", err),
		)
	}

	// wait for interrupt
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nClosing Peril server\nGoodbye!")
}
