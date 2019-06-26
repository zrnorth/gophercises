package deck

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Spade})
	fmt.Println(Card{Rank: Five, Suit: Diamond})
	fmt.Println(Card{Rank: Jack, Suit: Club})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Two of Spades
	// Five of Diamonds
	// Jack of Clubs
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	// 13 ranks, 4 suits
	if len(cards) != 13*4 {
		t.Error("Wrong number of cards in a new deck")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	expected := Card{Rank: Ace, Suit: Spade}
	if cards[0] != expected {
		t.Error("Expected default sorted deck to begin with Ace of Spaces. Received:", cards[0])
	}
}

func TestShuffle(t *testing.T) {
	shuffleRand = rand.New(rand.NewSource(0)) // Overwrote shuffleRand here to fix the seed
	// With key 0, order of cards is [40 35 ...]
	orig := New()
	first := orig[40]
	second := orig[35]

	cards := New(Shuffle)
	if cards[0] != first {
		t.Errorf("Expected the first shuffled card to be %s, received %s.", first, cards[0])
	}
	if cards[1] != second {
		t.Errorf("Expected the first shuffled card to be %s, received %s.", second, cards[1])
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(3))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}
	if count != 3 {
		t.Error("Expected 3 jokers, received:", count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}
	cards := New(Filter(filter))
	for _, c := range cards {
		if c.Rank == Two || c.Rank == Three {
			t.Error("Expected all Twos and Threes to be filtered out!")
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))
	if len(cards) != 13*4*3 {
		t.Errorf("Expected deck length of %d but received deck of %d cards", 13*4*3, len(cards))
	}
}
