package cli

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitInstanceFlag string
)

func init() {
	rabbitMQCmd.Flags().StringVarP(&rabbitInstanceFlag, "instance", "i", "client", "RabbitMQ instance: client or server")

	rootCmd.AddCommand(rabbitMQCmd)
}

var rabbitMQCmd = &cobra.Command{
	Use:   "rabbitmq",
	Short: "Start RabbitMQ",
	Long:  `Start RabbitMQ`,
	Run: func(cmd *cobra.Command, args []string) {
		if rabbitInstanceFlag == "server" {
			startRabbitMQServer()
		} else if rabbitInstanceFlag == "client" {
			startRabbitMQClient()
		} else {
			log.Fatalf("Invalid instance: %s (must be client or server)", rabbitInstanceFlag)
		}
	},
}

func startRabbitMQClient() {
	fmt.Printf("\nStart RabbitMQ client\n")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	var forever = make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func startRabbitMQServer() {
	fmt.Printf("\nStart RabbitMQ server\n")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
	log.Printf(" [x] Sent %s\n", body)
}
