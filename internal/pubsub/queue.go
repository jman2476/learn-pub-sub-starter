package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	queue, err := channel.QueueDeclare(
		queueName,                       // name
		queueType == SimpleQueueDurable, // durable
		queueType != SimpleQueueDurable, // autoDelete
		queueType != SimpleQueueDurable, // exclusive
		false,                           // noWait
		amqp.Table{
			"x-dead-letter-exchange": DeadLetterExchange, // args
		},
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = channel.QueueBind(
		queueName,
		key,
		exchange,
		false, nil,
	)

	return channel, queue, nil
}
