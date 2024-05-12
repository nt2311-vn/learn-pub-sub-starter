package pubsub

import (
	"log"

	amqb "github.com/rabbitmq/amqp091-go"
)

type SimpeQueueType int

const (
	SimpleQueueDurable SimpeQueueType = iota
	SimpleQueueTrasient
)

func DeclareAndBind(
	conn *amqb.Connection,
	exchange, queueName, key string,
	simpleQueueType SimpeQueueType,
) (*amqb.Channel, amqb.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		log.Printf("Cannot create channel: %v\n", err)
		return nil, amqb.Queue{}, err
	}

	queue, err := channel.QueueDeclare(
		queueName,
		simpleQueueType == SimpleQueueDurable,
		simpleQueueType != SimpleQueueDurable,
		simpleQueueType != SimpleQueueDurable,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Cannot create queue: %v\n", err)
		return nil, amqb.Queue{}, err
	}

	if err = channel.QueueBind(queueName, key, exchange, false, nil); err != nil {
		log.Printf("Cannot bind queue and channel: %v\n", err)
		return nil, amqb.Queue{}, nil
	}

	return channel, queue, nil
}
