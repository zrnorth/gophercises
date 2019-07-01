package blackjack

import (
	"errors"
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

	player    []playerHand
	handIdx   int
	playerBet int
	balance   int

	dealer   []deck.Card
	dealerAI AI
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player[g.handIdx].cards
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

type playerHand struct {
	cards []deck.Card
	bet   int
}

func bet(g *Game, ai AI, newlyShuffled bool) {
	bet := ai.Bet(newlyShuffled)
	if bet < 100 {
		panic("bet must be at least 100") // for testing
	}
	g.playerBet = bet
}

func deal(g *Game) {
	playerCards := make([]deck.Card, 0, 5)
	g.handIdx = 0
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.shoe = g.shoe[0], g.shoe[1:]
		playerCards = append(playerCards, card)
		card, g.shoe = g.shoe[0], g.shoe[1:]
		g.dealer = append(g.dealer, card)
	}
	g.player = []playerHand{
		{
			cards: playerCards,
			bet:   g.playerBet,
		},
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
			endRound(g, ai)
			continue
		}

		// Player's action
		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(*g.currentHand()))
			copy(hand, *g.currentHand())
			move := ai.Play(hand, g.dealer[0])
			err := move(g)

			switch err {
			case errBust:
				MoveStand(g)
			case nil:
				// noop
			default:
				panic(err)
			}
		}

		// Dealer's action
		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}

		endRound(g, ai)
	}
	return g.balance
}

var (
	errBust = errors.New("hand score exceeded 21")
)

type Move func(*Game) error

func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card
	card, g.shoe = draw(g.shoe)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		return errBust
	}
	return nil
}

func MoveStand(g *Game) error {
	if g.state == stateDealerTurn {
		g.state++
		return nil
	}
	if g.state == statePlayerTurn {
		g.handIdx++
		if g.handIdx == len(g.player) { // we have played all our hands
			g.state++
		}
		return nil
	}
	return errors.New("invalid state")
}

func MoveSplit(g *Game) error {
	cards := g.currentHand()
	if len(*cards) != 2 {
		return errors.New("can only split with exactly two of the same card")
	}
	if (*cards)[0].Rank != (*cards)[1].Rank {
		return errors.New("both cards must have the same rank to split")
	}
	g.player = append(g.player, playerHand{
		cards: []deck.Card{(*cards)[1]},
		bet:   g.player[g.handIdx].bet,
	})
	g.player[g.handIdx].cards = (*cards)[:1]
	return nil
}

func MoveDoubleDown(g *Game) error {
	if len(*g.currentHand()) != 2 {
		return errors.New("can only double down on a hand with exactly 2 cards")
	}
	g.playerBet *= 2
	MoveHit(g)
	return MoveStand(g)
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

func endRound(g *Game, ai AI) {
	dealerScore := Score(g.dealer...)
	dealerHasBlackjack := Blackjack(g.dealer...)
	allHands := make([][]deck.Card, len(g.player))
	for i, hand := range g.player {
		cards := hand.cards
		allHands[i] = cards

		winnings := hand.bet
		playerScore := Score(cards...)
		playerHasBlackjack := Blackjack(cards...)
		switch {
		case playerHasBlackjack && dealerHasBlackjack:
			winnings = 0
		case dealerHasBlackjack:
			winnings *= -1
		case playerHasBlackjack:
			winnings = int(float64(winnings) * g.blackjackPayout)
		case playerScore > 21:
			winnings *= -1
		case dealerScore > 21:
			// win
		case playerScore > dealerScore:
			// win
		case playerScore < dealerScore:
			winnings *= -1
		case playerScore == dealerScore:
			winnings = 0
		}
		g.balance += winnings
	}

	fmt.Println()
	ai.Results(allHands, g.dealer)
	g.player = nil
	g.dealer = nil
}
