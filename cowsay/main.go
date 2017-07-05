package main

import (
	"os"
	"log"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/dhruvbird/go-cowsay"
	"gopkg.in/matryer/try.v1"
)

var (
	connection *amqp.Connection
	channel *amqp.Channel
	task_q amqp.Queue
	response_q amqp.Queue
)


// from rabbitmqt/rabbitmq-tutorials example
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func init() {
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		connection, err = amqp.Dial(os.Getenv("AMQP_URL"))
		if err != nil {
			log.Println("Something happened during connection. Retry in 2 seconds...")
			time.Sleep(2 * time.Second)
		}
		return attempt < 5, err
	})
	if err != nil {
		fmt.Println("Connection error")
		log.Fatal(err)
	}
	log.Println("MQ connected.")

	channel, err = connection.Channel()
	if err != nil {
		fmt.Println("Connection error")
		log.Fatal(err)
	}

	task_q, err = channel.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	response_q, err = channel.QueueDeclare(
		"response_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = channel.Qos(
		1,
		0,
		false,
	)
	failOnError(err, "Failed to set QoS")
}

func main() {
	c, err := channel.Consume("task_queue", "", false, false, false, false, nil)
	failOnError(err, "Failed to consume message")

	for message := range c {
		messageId := message.Headers["MessageId"].(string)
		log.Printf("Consumer %q %q", messageId, message.Body)

		headers := amqp.Table{
			"MessageId": messageId,
		}
		cowsayMessage := cowsay.Format(string(message.Body))
		responseMessage := amqp.Publishing{
			Timestamp: time.Now(),
			ContentType: "text/plain",
			Body: []byte(cowsayMessage),
			Headers: headers,
		}
		channel.Publish("", "response_queue", false, false, responseMessage)
		message.Ack(false)
	}
	log.Printf("Ended.")
}
