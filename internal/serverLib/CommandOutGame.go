package serverLib

import (
	"encoding/json"
)

type CmdOutCardHint struct {
	Card []string `json:"card"`
}

func (cardHint *CmdOutCardHint) String() string {
	s, err := json.Marshal(cardHint)
	if err != nil {
		return "cardHint {}"
	}
	return "cardHint " + string(s)
}

type CmdOutYourTurn struct {
	Turn bool `json:"turn"`
}

func (c *CmdOutYourTurn) String() string {
	s, err := json.Marshal(c)
	if err != nil {
		return "yourTurn {}"
	}
	return "yourTurn " + string(s)
}

type CmdOutGameOver struct {
	CoverCard []string `json:"coverCard"`
	Score     []int    `json:"score"`
}

func (c *CmdOutGameOver) String() string {
	s, err := json.Marshal(c)
	if err != nil {
		return "gameOver {}"
	}
	return "gameOver " + string(s)
}

type CmdOutBackToRoom struct {
}

func (c *CmdOutBackToRoom) String() string {
	return "backToRoom"
}
