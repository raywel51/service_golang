package services

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"playground/configs"
	"time"
)

func ConnectToRabbitMQ(message string) error {
	err := configs.LoadEnv()
	if err != nil {
		return nil
	}

	brokerURL := "amqp://" + os.Getenv("MQTT_USERNAME") + ":" + os.Getenv("MQTT_PASSWORD") + "@" + os.Getenv("MQTT_HOSTNAME")
	queueName := os.Getenv("MQTT_QUEUE")

	if queueName == "" {
		return fmt.Errorf("MQTT_QUEUE environment variable is not defined")
	}

	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",        // Exchange
		queueName, // Routing key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("\nConnected to RabbitMQ\n")
	fmt.Printf("Sent Message to RabbitMQ: %s %s\n", message, time.Now().Format(time.RFC3339))

	return nil
}
