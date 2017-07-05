package main

import (
	"os"
	"net/http"
	"log"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"gopkg.in/matryer/try.v1"
)

var (
	connection *amqp.Connection
	channel *amqp.Channel
	task_q amqp.Queue
	response_q amqp.Queue

	responseChannels = make(map[string]chan string)
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


	go consumer()
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", r.FormValue("message"))
}

func cowsayHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("message")
	log.Printf("Received response: %q", msg)
	messageId := msg

	responseChannel := make(chan string)
	responseChannels[messageId] = responseChannel

	// :TODO: testing decoy message
	headers := amqp.Table{
		"MessageId": messageId,
	}
	message := amqp.Publishing{
		Timestamp: time.Now(),
		ContentType: "text/plain",
		Body: []byte(msg),
		Headers: headers,
	}
	channel.Publish("", "task_queue", false, false, message)

	response := <-responseChannel
	fmt.Fprintf(w, "%s", response)
}

func consumer() {
	c, err := channel.Consume("response_queue", "", false, false, false, false, nil)
	failOnError(err, "Failed to consume message")

	for message := range c {
		messageId := message.Headers["MessageId"].(string)
		log.Printf("Consumer %q %q", messageId, message.Body)

		responseChannel := responseChannels[messageId]
		if responseChannel != nil {
			log.Printf("Found a relevant channel")
			responseChannel <- string(message.Body)
			delete(responseChannels, messageId)
		} else {
			log.Printf("Failed to found a channel")
		}
		message.Ack(false)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/echo", echoHandler)
	mux.HandleFunc("/say", cowsayHandler)

	s := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Fatal(s.ListenAndServe())
}
