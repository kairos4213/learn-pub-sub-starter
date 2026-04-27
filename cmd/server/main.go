package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	connectionString = "amqp://guest:guest@localhost:5672/"
	exchange         = routing.ExchangePerilDirect
	pauseKey         = routing.PauseKey
)

func main() {
	fmt.Println("Starting Peril Server...")

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("server could not connect to RabbitMQ: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("server error closing connection to RabbitMQ: %v", err)
		}
	}()
	fmt.Println("Server Connection to RabbitMQ successful")

	pauseChann, err := conn.Channel()
	if err != nil {
		log.Printf("error opening pause/resume channel: %v", err)
	}
	if err = pubsub.PublishJSON(pauseChann, exchange, pauseKey, routing.PlayingState{IsPaused: true}); err != nil {
		log.Printf("error with pause channel: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("")
	fmt.Println("Peril Server successfully shutdown")
}
