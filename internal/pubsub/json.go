package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](
	ch *amqp.Channel,
	exchange,
	key string,
	val T,
) error {
	return publishData(
		ch,
		"application/json",
		exchange,
		key,
		val,
	)
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	return subscribeChannel(
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		unmarshalJSON,
	)
}
