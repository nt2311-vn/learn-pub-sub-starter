package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqb "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	connectionStr := "amqp://guest:guest@localhost:5672/"
	con, err := amqb.Dial(connectionStr)
	if err != nil {
		log.Fatalf("Cannot create channel connection: %v\n", err)
	}

	defer con.Close()
	fmt.Println("Peril game connected to rabbitmq")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Could not get username: %v\n", err)
	}

	_, que, err := pubsub.DeclareAndBind(
		con,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTrasient,
	)
	if err != nil {
		log.Fatalf("Could not subscribe to pause: %v\n", err)
	}

	fmt.Printf("Queue %v declared and bound\n", que.Name)

	state := gamelogic.NewGameState(username)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "move":
			_, err := state.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// TODO: Publish the move
		case "spawn":
			err = state.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			state.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			// TODO: publish n malicious logs
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}

	}
}
