package main

import (
	"fmt"

	blackjack "github.com/zrnorth/gopher/gopher-11/blackjack"
	deck "github.com/zrnorth/gopher/gopher-9"
)

type basicAI struct{}

func (ai basicAI) Bet(newlyShuffled bool) int {
	return 100
}

func (ai basicAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	if len(hand) == 2 { // Decide if we want to double or not
		if hand[0] == hand[1] { // Decide if we want to split or not
			cardScore := blackjack.Score(hand[0])
			if cardScore >= 8 && cardScore != 10 {
				return blackjack.MoveSplit
			}
		}
		if (score == 10 || score == 11) && !blackjack.Soft(hand...) {
			return blackjack.MoveDoubleDown
		}
	}
	// Decide if we want to hit or stand
	dealerScore := blackjack.Score(dealer)
	if dealerScore >= 5 && dealerScore <= 6 {
		return blackjack.MoveStand
	}
	if score < 13 {
		return blackjack.MoveHit
	}
	return blackjack.MoveStand
}

func (ai basicAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	// noop for now
}

func main() {
	opts := blackjack.Options{
		Decks:           4,
		Hands:           50000,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(opts)
	winnings := game.Play(&basicAI{})
	fmt.Println(winnings)
}
