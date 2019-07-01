package main

import (
	"fmt"

	blackjack "github.com/zrnorth/gopher/gopher-11/blackjack"
	deck "github.com/zrnorth/gopher/gopher-9"
)

type basicAI struct {
	count int
	seen  int
	decks int
}

func (ai *basicAI) Bet(newlyShuffled bool) int {
	if newlyShuffled {
		ai.count = 0
		ai.seen = 0
	}
	trueScore := ai.count / (((ai.decks * 52) - ai.seen) / 52)
	switch {
	case trueScore >= 14:
		return 1000
	case trueScore >= 8:
		return 500
	default:
		return 100
	}
}

func (ai *basicAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
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

func (ai *basicAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	for _, card := range dealer {
		ai.countCard(card)
	}
	for _, hand := range hands {
		for _, card := range hand {
			ai.countCard(card)
		}
	}
}

func (ai *basicAI) countCard(card deck.Card) {
	score := blackjack.Score(card)
	// Basic count
	switch {
	case score >= 10:
		ai.count--
	case score <= 6:
		ai.count++
	}

	ai.seen++
}

func main() {
	opts := blackjack.Options{
		Decks:           4,
		Hands:           50000,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(opts)
	winnings := game.Play(&basicAI{
		decks: opts.Decks,
	})
	fmt.Println(winnings)
}
