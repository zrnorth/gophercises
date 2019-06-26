package blackjack

import (
	"fmt"

	deck "github.com/zrnorth/gopher/gopher-9"
)

type AI interface {
	Bet() int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hand [][]deck.Card, dealer []deck.Card)
}

func HumanAI() AI {
	return humanAI{}
}

type humanAI struct{}

type dealerAI struct {
}

func (ai dealerAI) Bet() int {
	return 1 // noop, dealer doesn't bet
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	} else {
		return MoveStand
	}
}

func (ai dealerAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	// noop
}

func (ai humanAI) Bet() int {
	return 1 // todo
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	// dealer deck.Card is not a slice bc when player is making move
	// they can only see 1 dealer card
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)
		fmt.Println("What will you do? (h)it, (s)tand")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		default:
			fmt.Println("Invalid option:", input)
		}
	}
}

func (ai humanAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	// Player can have multiple hands because of splitting
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", hand)
	fmt.Println("Dealer:", dealer)
}
