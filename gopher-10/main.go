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

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Shoe   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Shoe:   make([]deck.Card, len(gs.Shoe)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(ret.Shoe, gs.Shoe)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)
	return ret
}

func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Shoe = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Shoe = ret.Shoe[0], ret.Shoe[1:]
		ret.Player = append(ret.Player, card)
		card, ret.Shoe = ret.Shoe[0], ret.Shoe[1:]
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.State = StatePlayerTurn
	return ret
}

func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()

	var card deck.Card
	card, ret.Shoe = ret.Shoe[0], ret.Shoe[1:]
	*hand = append(*hand, card)

	// Could early-out here if the player is now over 21.
	return ret
}

func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++ // Additional states could break this
	return ret
}

func EndHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := gs.Player.Score(), gs.Dealer.Score()
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", gs.Player, "\nScore:", pScore)
	fmt.Println("Dealer:", gs.Dealer, "\nScore:", dScore)
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
	fmt.Println()
	ret.Player = nil
	ret.Dealer = nil
	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)
	gs = Deal(gs)

	// Player's action
	var input string
	for gs.State == StatePlayerTurn {
		fmt.Println("Player:", gs.Player)
		fmt.Println("Dealer:", gs.Dealer.DealerString())
		fmt.Println("What will you do? (h)it, (s)tand")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			gs = Hit(gs)
		case "s":
			gs = Stand(gs)
		default:
			fmt.Println("Invalid option:", input)
		}
	}

	// Dealer's action
	for gs.State == StateDealerTurn {
		if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}
	}

	gs = EndHand(gs)
}
