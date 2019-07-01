package main

import (
	"fmt"

	blackjack "github.com/zrnorth/gopher/gopher-11/blackjack"
)

func main() {
	game := blackjack.New(blackjack.Options{})
	winnings := game.Play(blackjack.HumanAI())
	fmt.Println(winnings)
}
