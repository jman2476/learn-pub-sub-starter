package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func publishData[T any](
	ch *amqp.Channel,
	contentType,
	exchange,
	key string,
	val T,
) error {
	var data []byte
	var err error
	switch contentType {
	case "application/json":
		data, err = encodeJSON(val)
		if err != nil {
			return err
		}
	case "application/gob":
		data, err = encodeGob(val)
		if err != nil {
			return nil
		}
	}

	return ch.PublishWithContext(
		context.Background(),
		exchange, key,
		false, false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        data,
		},
	)
}

func encodeJSON[T any](val T) ([]byte, error) {
	return json.Marshal(val)
}

func encodeGob[T any](val T) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(val)
	if err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}
