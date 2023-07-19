package serverLib

import (
	"encoding/json"
	"strings"
	"sync"
)

type Card struct {
	hands map[string]int
	Cover []string
	mutex sync.RWMutex
}

func NewCard(hands map[string]int) *Card {
	return &Card{
		hands: hands,
		Cover: []string{},
	}
}

func (c *Card) IsEmpty() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.hands) == 0
}

func (c *Card) GetCard() map[string]int {
	return c.hands
}

func (c *Card) CoverCard(card string) {
	c.Cover = append(c.Cover, card)
	delete(c.hands, card)
}

func (c *Card) RemoveCard(card string) {
	delete(c.hands, card)
}

func (c *Card) String() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	s := ""
	for k, _ := range c.hands {
		s += k
	}
	data := struct {
		Card      string `json:"card"`
		CoverCard string `json:"coverCard"`
	}{
		Card:      s,
		CoverCard: strings.Join(c.Cover[:], ""),
	}
	b, _ := json.Marshal(data)

	return "cardInfo " + string(b)
}
