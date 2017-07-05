package main

import (
	"os"
	"net/http"
	"log"
	"fmt"
	"html"

	"github.com/streadway/amqp"
)

var (
	connection *amqp.Connection
	channel *amqp.Channel
)

func init() {
	connection, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		fmt.Println("Connection error")
		log.Fatal(err)
	}

	channel, err = connection.Channel()
	if err != nil {
		fmt.Println("Connection error")
		log.Fatal(err)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", r.FormValue("message"))
}

func cowsayHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
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
