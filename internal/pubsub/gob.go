package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](
	ch *amqp.Channel,
	exchange,
	key string,
	val T,
) error {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.Background(),
		exchange, key,
		false, false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        data.Bytes(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
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

		for msg := range delchan {
			data, err := unmarshaller([]byte(msg.Body))
			if err != nil {
				fmt.Println(
					fmt.Errorf("Error unmarshaling: %w", err),
				)
			}
			ackType := handler(data)
			switch ackType {
			case Ack:
				msg.Ack(false)
			case NackRequeue:
				msg.Nack(false, true)
			case NackDiscard:
				msg.Nack(false, false)
			}
		}
	}(deliveryChan)

	return nil
}

func UnmarshalGob[T any](data []byte) (T, error) {
	var buf = bytes.NewBuffer(data)
	var t T
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&t)
	return t, err
}
