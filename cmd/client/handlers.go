package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jman2476/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jman2476/learn-pub-sub-starter/internal/pubsub"
	"github.com/jman2476/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(mv gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print(">")
		outcome := gs.HandleMove(mv)

		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				strings.Join([]string{routing.WarRecognitionsPrefix, gs.Player.Username}, "."),
				gamelogic.RecognitionOfWar{
					Attacker: mv.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			fallthrough
		default:
			return pubsub.NackDiscard
		}
	}
}

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(row gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(row gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")

		outcome, winner, loser := gs.HandleWar(row)
		var logMsg string
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			fallthrough
		case gamelogic.WarOutcomeYouWon:
			logMsg = fmt.Sprintf("%s won the war against %s", winner, loser)
			fallthrough
		case gamelogic.WarOutcomeDraw:
			if logMsg == "" {
				logMsg = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			}
			err := PublishLogs(
				routing.GameLog{
					CurrentTime: time.Now(),
					Message:     logMsg,
					Username:    gs.Player.Username,
				},
				ch,
			)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			fmt.Print(
				fmt.Errorf("Something went wrong with the war. Maybe aliens invaded?"),
			)
			return pubsub.NackDiscard
		}
	}

}
