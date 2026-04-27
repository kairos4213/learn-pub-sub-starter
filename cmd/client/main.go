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

	gameState := gamelogic.NewGameState(username)
	for {
		input := gamelogic.GetInput()
		switch input[0] {
		case "spawn":
			err = gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			_, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Invalid Command")
		}
	}
}
