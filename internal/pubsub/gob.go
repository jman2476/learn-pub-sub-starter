package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](
	ch *amqp.Channel,
	exchange,
	key string,
	val T,
) error {
	return publishData(
		ch,
		"application/gob",
		exchange,
		key,
		val,
	)
}

func SubscribeGob[T any](
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
		unmarshalGob,
	)
}
