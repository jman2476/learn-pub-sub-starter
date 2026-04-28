package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](
	ch *amqp.Channel,
	exchange,
	key string,
	val T,
) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.Background(),
		exchange, key,
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {
	channel, _, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		return err
	}

	deliveryChan, err := channel.Consume(
		"", "",
		false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	go func(delchan <-chan amqp.Delivery) {
		for del := range delchan {
			var data T
			err := json.Unmarshal([]byte(del.Body), &data)
			if err != nil {
				fmt.Println(
					fmt.Errorf("Error unmarshaling: %w", err),
				)
			}
			handler(data)
			del.Ack(false)
		}
	}(deliveryChan)

	return nil
}
