package serverLib

import "sync"

type GameDeck struct {
	card  map[string]string
	mutex sync.RWMutex
}

func NewGameDeck() *GameDeck {
	card := map[string]string{"1": "", "2": "", "3": "", "4": ""}
	return &GameDeck{
		card: card,
	}
}
