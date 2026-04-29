package main

import (
	"fmt"

	"github.com/jman2476/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jman2476/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(mv gamelogic.ArmyMove) {
		defer fmt.Print(">")
		gs.HandleMove(mv)
	}
}
