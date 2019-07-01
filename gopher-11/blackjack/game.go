package blackjack

import (
	"fmt"

	deck "github.com/zrnorth/gopher/gopher-9"
)

type state int8

type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

const (
	stateBetting state = iota
	statePlayerTurn
	stateDealerTurn
	stateHandOver
)

func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}

	// Set default options
	if opts.Decks == 0 {
		opts.Decks = 3
	}
	if opts.Hands == 0 {
		opts.Hands = 100
	}
	if opts.BlackjackPayout == 0 {
		opts.BlackjackPayout = 1.5
	}

	g.numDecks = opts.Decks
	g.numHands = opts.Hands
	g.blackjackPayout = opts.BlackjackPayout

	return g
}

type Game struct {
	// unexported fields
	numDecks        int // number of decks in the shoe
	numHands        int
	blackjackPayout float64

	state state
	shoe  []deck.Card

	player    []deck.Card
	playerBet int
	balance   int

	dealer   []deck.Card
	dealerAI AI
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func bet(g *Game, ai AI, newlyShuffled bool) {
	bet := ai.Bet(newlyShuffled)
	g.playerBet = bet
}

func deal(g *Game) {
	g.player = make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.shoe = g.shoe[0], g.shoe[1:]
		g.player = append(g.player, card)
		card, g.shoe = g.shoe[0], g.shoe[1:]
		g.dealer = append(g.dealer, card)
	}
	g.state = statePlayerTurn
}

func (g *Game) Play(ai AI) int {
	g.shoe = nil
	min := (52 * g.numDecks) / 2 // when shoe size hits min, refshuffle
	for i := 0; i < g.numHands; i++ {
		newlyShuffled := false
		if len(g.shoe) < min {
			g.shoe = deck.New(deck.Deck(g.numDecks), deck.Shuffle)
			newlyShuffled = true
		}
		bet(g, ai, newlyShuffled)
		deal(g)
		// Check for dealer blackjack and early out if so
		if Blackjack(g.dealer...) {
			endHand(g, ai)
			continue
		}

		// Player's action
		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			move(g)
		}

		// Dealer's action
		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}

		endHand(g, ai)
	}
	return g.balance
}

type Move func(*Game)

func MoveHit(g *Game) {
	hand := g.currentHand()
	var card deck.Card
	card, g.shoe = draw(g.shoe)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		MoveStand(g)
	}
	return
}

func MoveStand(g *Game) {
	g.state++ // todo this will change
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// Soft returns true if a blackjack hand contains an A being counted as 11.
func Soft(hand ...deck.Card) bool {
	return minScore(hand...) != Score(hand...)
}

// Blackjack returns true if a hand is a blackjack (Ace + 10 card)
func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

// Score returns the best possible blackjack hand for a given slice of cards.
func Score(hand ...deck.Card) int {
	minSc := minScore(hand...)
	if minSc > 11 {
		return minSc
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			return minSc + 10 // A goes from 1->11
		}
	}
	return minSc
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, c := range hand {
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

func endHand(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)
	playerHasBlackjack, dealerHasBlackjack := Blackjack(g.player...), Blackjack(g.dealer...)
	winnings := g.playerBet
	switch {
	case playerHasBlackjack && dealerHasBlackjack:
		winnings = 0
	case dealerHasBlackjack:
		winnings *= -1
	case playerHasBlackjack:
		winnings = int(float64(winnings) * g.blackjackPayout)
	case pScore > 21:
		winnings *= -1
	case dScore > 21:
		// win
	case pScore > dScore:
		// win
	case pScore < dScore:
		winnings *= -1
	case pScore == dScore:
		winnings = 0
	}
	g.balance += winnings

	fmt.Println()
	ai.Results([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.dealer = nil
}
