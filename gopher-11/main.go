package main

import (
	"fmt"

	blackjack "github.com/zrnorth/gopher/gopher-11/blackjack"
)

func main() {
	opts := blackjack.Options{
		Decks:           2,
		Hands:           2,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(opts)
	winnings := game.Play(blackjack.HumanAI())
	fmt.Println(winnings)
}
