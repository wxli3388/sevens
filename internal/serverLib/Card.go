package serverLib

import "encoding/json"

type Card struct {
	hands map[string]int
	Cover map[int]int
}

func NewCard(hands map[string]int) *Card {
	return &Card{
		hands: hands,
		Cover: make(map[int]int),
	}
}

func (c *Card) GetCard() map[string]int {
	return c.hands
}

func (c *Card) String() string {
	s := ""
	for k, _ := range c.hands {
		s += k
	}
	data := struct {
		Card string `json:"card"`
	}{
		Card: s,
	}
	b, _ := json.Marshal(data)

	return "cardInfo " + string(b)
}
