package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	connectionString = "amqp://guest:guest@localhost:5672/"
	directExchange   = routing.ExchangePerilDirect
	topicExchange    = routing.ExchangePerilTopic
	pauseKey         = routing.PauseKey
	logKey           = routing.GameLogSlug
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

	_, _, err = pubsub.DeclareAndBind(conn, topicExchange, "game_logs", logKey+".*", pubsub.Durable)
	if err != nil {
		log.Printf("error declaring and binding game logs: %v", err)
	}

	pauseChann, err := conn.Channel()
	if err != nil {
		log.Printf("error opening pause/resume channel: %v", err)
	}

	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		switch input[0] {
		case "pause":
			log.Println("Sending pause message")
			if err = pubsub.PublishJSON(pauseChann, directExchange, pauseKey, routing.PlayingState{IsPaused: true}); err != nil {
				log.Printf("error publishing pause message to pause channel: %v", err)
			}
		case "resume":
			log.Println("Sending resume message")
			if err = pubsub.PublishJSON(pauseChann, directExchange, pauseKey, routing.PlayingState{IsPaused: false}); err != nil {
				log.Printf("error publishing resume message to pause channel: %v", err)
			}
		case "quit":
			log.Println("Sending exit message")
			return
		default:
			log.Println("Invalid command")
		}
	}
}
