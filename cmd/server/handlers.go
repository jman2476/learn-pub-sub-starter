package main

import (
	"fmt"

	"github.com/jman2476/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jman2476/learn-pub-sub-starter/internal/pubsub"
	"github.com/jman2476/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerLogs(logQueue amqp.Queue) func(routing.GameLog) pubsub.Acktype {
	return func(gl routing.GameLog) pubsub.Acktype {
		defer fmt.Print("> ")
		fmt.Printf("Logging on %s\n", logQueue.Name)
		err := gamelogic.WriteLog(gl)
		if err != nil {
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
