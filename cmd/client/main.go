package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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

	channel, queue, err := pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		strings.Join([]string{routing.PauseKey, username}, "."),
		routing.PauseKey,
		"transient",
	)
	if err != nil {
		fmt.Println(
			fmt.Errorf("Error getting channel and queue: %w", err),
		)
	}

	fmt.Printf("Channel: %v Queue: %v", channel, queue)
	// wait for interrupt
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nClosing Peril client\nGoodbye!")
}
