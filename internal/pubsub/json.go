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

	// data, err := json.Marshal(val)
	// if err != nil {
	// 	return err
	// }

	// err = ch.PublishWithContext(
	// 	context.Background(),
	// 	exchange, key,
	// 	false, false,
	// 	amqp.Publishing{
	// 		ContentType: "application/json",
	// 		Body:        data,
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }

	// return nil
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

	// channel, _, err := DeclareAndBind(
	// 	conn,
	// 	exchange,
	// 	queueName,
	// 	key,
	// 	queueType,
	// )
	// if err != nil {
	// 	return err
	// }

	// deliveryChan, err := channel.Consume(
	// 	"", "",
	// 	false, false, false, false, nil,
	// )
	// if err != nil {
	// 	return err
	// }

	// go func(delchan <-chan amqp.Delivery) {

	// 	for msg := range delchan {
	// 		var data T
	// 		err := json.Unmarshal([]byte(msg.Body), &data)
	// 		if err != nil {
	// 			fmt.Println(
	// 				fmt.Errorf("Error unmarshaling: %w", err),
	// 			)
	// 		}
	// 		ackType := handler(data)
	// 		switch ackType {
	// 		case Ack:
	// 			msg.Ack(false)
	// 		case NackRequeue:
	// 			msg.Nack(false, true)
	// 		case NackDiscard:
	// 			msg.Nack(false, false)
	// 		}
	// 	}
	// }(deliveryChan)

	// return nil
}
