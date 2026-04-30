package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribeChannel[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshamller func([]byte) (T, error),
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

	channel.Qos(10, 0, true)
	deliveryChan, err := channel.Consume(
		"", "",
		false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	go func(delchan <-chan amqp.Delivery) {

		for msg := range delchan {
			data, err := unmarshamller([]byte(msg.Body))
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

func unmarshalGob[T any](data []byte) (T, error) {
	var buf = bytes.NewBuffer(data)
	var t T
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&t)
	return t, err
}

func unmarshalJSON[T any](data []byte) (T, error) {
	var processed T
	err := json.Unmarshal(data, processed)

	return processed, err
}
