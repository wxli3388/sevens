package serverLib

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

type Player struct {
	isRobot    bool
	User       *User
	Card       *Card
	position   int
	game       *Game
	gameAction chan *GameAction
	endGame    chan struct{}
}

func NewPlayer(user *User) *Player {
	isRobot := false
	if user == nil {
		isRobot = true
	}
	player := &Player{
		isRobot:    isRobot,
		User:       user,
		gameAction: make(chan *GameAction),
		endGame:    make(chan struct{}),
	}
	if !player.isRobot {
		go player.ServeCmd()
	}
	return player
}

const (
	RobotAutoPlay = 20    //ms
	HumanAutoPlay = 20000 //ms
)

func (p *Player) ServeCmd() {
	run := true
	for run {
		select {
		case cmd := <-p.User.gameCmd:
			p.handleGameCmd(cmd)
		case <-p.endGame:
			run = false
			break
		}
	}
}

func (p *Player) handleGameCmd(cmd string) {
	res := strings.SplitN(cmd, " ", 2)
	switch res[0] {
	case "game_cover_card":
		p.cmdGameCoverCard(res[1])
	case "game_play_card":
		p.cmdGamePlayCard(res[1])
	}
}

func (p *Player) cmdGameCoverCard(data string) {
	if !p.game.CheckTurn(p.position) {
		return // not your turn
	}
	var cmdInGameCoverCard CmdInGameCoverCard
	err := json.Unmarshal([]byte(data), &cmdInGameCoverCard)
	if err != nil {
		return
	}
	playCard := cmdInGameCoverCard.Card
	if !p.cardCheck(playCard) {
		return
	}

	ga := &GameAction{
		actionCode: GameActionCover,
		card:       playCard,
	}
	timer := time.NewTimer(500 * time.Millisecond) //防呆 應該能快速處理完?
	select {
	case p.gameAction <- ga:
		break
	case <-timer.C:
		break
	}
}

func (p *Player) cmdGamePlayCard(data string) {
	if !p.game.CheckTurn(p.position) {
		return // not your turn
	}
	var cmdInGamePlayCard CmdInGamePlayCard
	err := json.Unmarshal([]byte(data), &cmdInGamePlayCard)
	if err != nil {
		return
	}
	playCard := cmdInGamePlayCard.Card
	if !p.cardCheck(playCard) {
		return
	}

	valid := false
	canPlay := p.GetCardCanPlay()
	for _, v := range canPlay {
		if v == playCard {
			valid = true
			break
		}
	}
	if !valid {
		p.User.Write("You can't play this card!!")
		return
	}

	ga := &GameAction{
		actionCode: GameActionPlay,
		card:       playCard,
	}
	timer := time.NewTimer(500 * time.Millisecond) //防呆 應該能快速處理完?
	select {
	case p.gameAction <- ga:
		break
	case <-timer.C:
		break
	}
}

func (p *Player) cardCheck(card string) bool {

	if _, ok := p.Card.GetCard()[card]; !ok {
		p.User.Write("You don't have this card!!")
		return false
	}
	return true
}

func (p *Player) SetGame(game *Game) {
	p.game = game
}

func (p *Player) SetCard(card *Card) {
	p.Card = card
}

func (p *Player) SetPosition(position int) {
	p.position = position
}

func (p *Player) GetPosition() int {
	return p.position
}

func (p *Player) RoundStart() {
	if p.Card.IsEmpty() {
		p.game.finish <- struct{}{}
		return
	}
	if p.isRobot {
		timer := time.NewTimer(RobotAutoPlay * time.Millisecond)
		<-timer.C
		p.AutoPlay()
		p.turnNextPlayer()
		return
	}
	IntervalTime := HumanAutoPlay * time.Millisecond // 觸發間隔時間
	ticker := time.NewTicker(IntervalTime)           // 設定 秒觸發一次
	card := p.GetCardCanPlay()
	p.User.Write(&CmdOutYourTurn{true})
	p.User.Write(&CmdOutCardHint{card})
	run := true
	for run {
		select {
		case gameAction := <-p.gameAction:
			if gameAction.actionCode == GameActionPlay {
				p.PlayCard(gameAction.card)
			} else {
				p.CoverCard(gameAction.card)
			}

			p.game.Action <- gameAction
			p.UpdateUserCard()
			run = false
			break
		case <-ticker.C:
			p.AutoPlay()
			run = false
			break
		}
	}
	p.User.Write(&CmdOutYourTurn{false})
	p.turnNextPlayer()
}

func (p *Player) turnNextPlayer() {
	p.game.mutex.Lock()
	p.game.turn = (p.game.turn + 1) % p.game.maxPlayer
	p.game.mutex.Unlock()
	p.game.Players[p.game.turn].RoundStart() // remove goroutine, change turn with channel
}

func (p *Player) AutoPlay() {
	if _, ok := p.Card.GetCard()["47"]; ok {
		p.PlayCard("47")
		p.game.Action <- &GameAction{
			actionCode: GameActionPlay,
			card:       "47",
		}
		p.UpdateUserCard()
		return
	}
	canPlay := p.GetCardCanPlay()
	if len(canPlay) > 0 {
		p.PlayCard(canPlay[0])
		p.game.Action <- &GameAction{
			actionCode: GameActionPlay,
			card:       canPlay[0],
		}
		p.UpdateUserCard()
		return
	}

	p.Card.mutex.Lock()
	cardMap := p.Card.GetCard()
	cardList := make([]string, 0, len(cardMap))
	for k, _ := range cardMap {
		cardList = append(cardList, k)
	}
	min := 0
	max := len(cardList) - 1

	idx := rand.Intn(max-min+1) + min

	p.Card.CoverCard(cardList[idx])
	p.game.Action <- &GameAction{
		actionCode: GameActionCover,
		card:       cardList[idx],
	}
	p.Card.mutex.Unlock()
	p.UpdateUserCard()
}

func (p *Player) PlayCard(card string) {
	p.game.GameDeck.mutex.Lock()

	suits, point := string(card[0]), string(card[1])
	if p.game.GameDeck.card[suits] == "" {
		p.game.GameDeck.card[suits] += point
	} else if string(p.game.GameDeck.card[suits][0]) < point {
		p.game.GameDeck.card[suits] += point
	} else {
		p.game.GameDeck.card[suits] = point + p.game.GameDeck.card[suits]
	}
	p.game.GameDeck.mutex.Unlock()
	p.Card.RemoveCard(card)
}

func (p *Player) CoverCard(card string) {
	p.Card.CoverCard(card)
}

func (p *Player) GetCardCanPlay() []string {
	defer p.Card.mutex.Unlock()
	p.Card.mutex.Lock()
	card := p.Card.GetCard()
	if _, ok := card["47"]; ok {
		return []string{"47"}
	}

	canPlay := []string{}
	for suit, cardStr := range p.game.GameDeck.card {
		if cardStr == "" {
			seven := suit + "7"
			if _, ok := card[seven]; ok {
				canPlay = append(canPlay, seven)
			}
		} else {
			pre, end := string(cardStr[0]), string(cardStr[len(cardStr)-1])

			front, back := p.GetLess(pre), p.GetLarge(end)
			cardFront := suit + front
			cardBack := suit + back
			if _, ok := card[cardFront]; ok {
				canPlay = append(canPlay, cardFront)
			}
			if _, ok := card[cardBack]; ok {
				canPlay = append(canPlay, cardBack)
			}
		}
	}
	return canPlay
}

func (p *Player) UpdateUserCard() {
	if p.isRobot {
		return
	}
	p.User.Write(p.Card)
}

func (p *Player) GetLarge(point string) string {
	h := map[string]string{
		"7": "8",
		"8": "9",
		"9": "A",
		"A": "B",
		"B": "C",
		"C": "D",
	}
	return h[point]
}

func (p *Player) GetLess(point string) string {
	h := map[string]string{
		"7": "6",
		"6": "5",
		"5": "4",
		"4": "3",
		"3": "2",
		"2": "1",
	}
	return h[point]
}
