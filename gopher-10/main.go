package main

import (
	"fmt"
	"strings"

	deck "github.com/zrnorth/gopher/gopher-9"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func draw(shoe []deck.Card) (deck.Card, []deck.Card) {
	return shoe[0], shoe[1:]
}

func main() {
	// Create a shuffled deck
	shoe := deck.New(deck.Deck(3), deck.Shuffle)
	var card deck.Card
	var player, dealer Hand
	// Deal the cards to the player and the dealer
	for i := 0; i < 2; i++ {
		for _, hand := range []*Hand{&player, &dealer} { // iterate over player and dealer's hands
			card, shoe = draw(shoe)
			*hand = append(*hand, card)
		}
	}
	fmt.Println("Player:", player)
	fmt.Println("Dealer:", dealer)
}
