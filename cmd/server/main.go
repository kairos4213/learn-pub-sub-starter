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

func main() {
	const connectionString = "amqp://guest:guest@localhost:5672/"
	const exchange = routing.ExchangePerilDirect
	const pauseKey = routing.PauseKey
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection to RabbitMQ: %v", err)
		}
	}()
	fmt.Println("Connection to RabbitMQ successful")

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
	fmt.Println("Peril server shutdown")
}
