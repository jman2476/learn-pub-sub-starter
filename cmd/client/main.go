package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

}
