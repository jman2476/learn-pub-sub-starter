package main

import (
	"strings"

	"github.com/jman2476/learn-pub-sub-starter/internal/pubsub"
	"github.com/jman2476/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishLogs(gl routing.GameLog, ch *amqp.Channel) error {
	return pubsub.PublishGob(
		ch,
		routing.ExchangePerilTopic,
		strings.Join([]string{routing.GameLogSlug, gl.Username}, "."),
		gl,
	)
}
