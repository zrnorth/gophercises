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

// In blackjack the dealer only shows one card
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			return minScore + 10 // A goes from 1->11
		}
	}
	return minScore
}

func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10) // cards can not be worth more than 10 in blackjack
	}
	return score
}

// helper function to get the min of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

	// Player action
	var input string
	for input != "s" { // While the player has not yet chosen to Stand
		fmt.Println("Player:", player)
		fmt.Println("Dealer:", dealer.DealerString())
		fmt.Println("What will you do? (h)it, (s)tand")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			card, shoe = draw(shoe)
			player = append(player, card)
		}
	}

	// Dealer action
	for dealer.Score() <= 16 || (dealer.Score() == 17 && dealer.MinScore() != 17) {
		card, shoe = draw(shoe)
		dealer = append(dealer, card)
		fmt.Println("Dealer hits:", card)
	}

	pScore, dScore := player.Score(), dealer.Score()
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", player, "\nScore:", pScore)
	fmt.Println("Dealer:", dealer, "\nScore:", dScore)
	switch {
	case pScore > 21:
		fmt.Println("You busted.")
	case dScore > 21:
		fmt.Println("The dealer busted!")
	case pScore > dScore:
		fmt.Println("You win!")
	case pScore < dScore:
		fmt.Println("You lose.")
	case pScore == dScore:
		fmt.Println("Draw!")
	}
}
