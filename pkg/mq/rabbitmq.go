package mq

import (
	"github.com/streadway/amqp"
	"log"
)

// RabbitMQ struct to hold connection and channel
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ instance
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create a channel: %v", err)
		return nil, err
	}

	return &RabbitMQ{Connection: conn, Channel: ch}, nil
}

// DeclareQueue declares a queue
func (r *RabbitMQ) DeclareQueue(queueName string) (amqp.Queue, error) {
	queue, err := r.Channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
		return amqp.Queue{}, err
	}
	return queue, nil
}

// PublishMessage publishes a message to a queue
func (r *RabbitMQ) PublishMessage(queueName, message string) error {
	err := r.Channel.Publish(
		"",        // Exchange
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		return err
	}
	log.Printf("Sent: %s", message)
	return nil
}

// ConsumeMessages consumes messages from a queue
func (r *RabbitMQ) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	messages, err := r.Channel.Consume(
		queueName, // Queue name
		"",        // Consumer tag
		true,      // Auto-acknowledge
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
		return nil, err
	}
	return messages, nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() {
	err := r.Channel.Close()
	if err != nil {
		log.Fatalf("Failed to close channel: %v", err)
	}
	err = r.Connection.Close()
	if err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}
