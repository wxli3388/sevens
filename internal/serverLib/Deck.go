package serverLib

import (
	"math/rand"
	"time"
)

type Deck struct {
	deck []string
}

func CreateDeck() *Deck {
	deck := []string{}
	suits := []string{"1", "2", "3", "4"}
	values := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D"}

	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, suit+value)
		}
	}
	return &Deck{
		deck: deck,
	}
}

func (d *Deck) ShuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.deck), func(i, j int) {
		d.deck[i], d.deck[j] = d.deck[j], d.deck[i]
	})
}

func (d *Deck) DealCards(players []*Player) {
	numCards := len(d.deck) / len(players)

	for _, player := range players {
		cardList := d.deck[:numCards]
		cardMap := map[string]int{}
		for _, v := range cardList {
			cardMap[v] = 1
		}
		card := NewCard(cardMap)
		player.SetCard(card)
		d.deck = d.deck[numCards:]
	}
}
