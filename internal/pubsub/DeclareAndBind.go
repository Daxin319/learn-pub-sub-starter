package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

type SimpleQueueType string

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	newChannel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	var queue amqp.Queue

	if queueType == "durable" {
		queue, err = newChannel.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
	} else if queueType == "transient" {
		queue, err = newChannel.QueueDeclare(queueName, false, true, true, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
	}

	err = newChannel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return newChannel, queue, nil
}
