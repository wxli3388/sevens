package serverLib

import (
	"fmt"
	"time"
)

type Player struct {
	isRobot  bool
	User     *User
	Card     *Card
	position int
	game     *Game
}

func NewPlayer(user *User) *Player {
	isRobot := false
	if user == nil {
		isRobot = true
	}
	return &Player{
		isRobot: isRobot,
		User:    user,
	}
}

func (p *Player) SetCard(card *Card) {
	p.Card = card
}

func (p *Player) SetPosition(position int) {
	p.position = position
}

// func (p *Player) Run() {
// 	for {
// 	case cmd := <-p.User.gameCmd:

// 	}
// }
func (p *Player) RoundStart() {
	ch := make(chan string)
	select {
	case <-ch:
	case <-time.After(time.Second * 5):
		fmt.Println("timeout")

	}
}
