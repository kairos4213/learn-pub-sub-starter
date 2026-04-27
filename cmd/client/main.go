package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
	fmt.Println("Starting Peril Client...")

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("client could not connect to RabbitMQ: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("client error closing connection to RabbitMQ: %v", err)
		}
	}()
	fmt.Println("Client Connection to RabbitMQ successful")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Printf("error getting username: %v", err)
	}

	_, _, err = pubsub.DeclareAndBind(conn, exchange, pauseKey+"."+username, pauseKey, pubsub.Transient)
	if err != nil {
		log.Printf("error declaring and binding pause channel: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("")
	fmt.Println("Peril Client successfully shutdown")
}
