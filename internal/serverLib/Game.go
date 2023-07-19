package serverLib

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Game struct {
	Players   []*Player
	GameDeck  *GameDeck
	Action    chan *GameAction
	finish    chan struct{}
	turn      int
	maxPlayer int
	mutex     sync.RWMutex
}

const (
	GameActionPlay  = 1
	GameActionCover = 2
)

type GameAction struct {
	actionCode int
	card       string
}

func NewGame(userMap map[*User]bool, maxPlayer int) *Game {
	var players []*Player
	for user, _ := range userMap {
		players = append(players, NewPlayer(user))
	}
	for i := len(players); i < maxPlayer; i += 1 {
		players = append(players, NewPlayer(nil))
	}
	game := &Game{
		Players:   players,
		Action:    make(chan *GameAction),
		finish:    make(chan struct{}),
		GameDeck:  NewGameDeck(),
		maxPlayer: maxPlayer,
		turn:      0,
	}
	for i := 0; i < maxPlayer; i += 1 {
		game.Players[i].SetGame(game)
	}
	return game
}

func (game *Game) Start() {
	deck := CreateDeck()
	deck.ShuffleDeck()
	deck.DealCards(game.Players)
	game.RandomPlayer()
	game.Broadcast("gameStart")
	game.SendCardStatus()
	go game.Run()

	go game.Players[game.turn].RoundStart()
}

func (game *Game) Run() {
	run := true
	for run {
		select {
		case action := <-game.Action:
			if action.actionCode == GameActionPlay {
				game.GameDeck.mutex.Lock()
				data := struct {
					Card     string            `json:"play_card"`
					DeskCard map[string]string `json:"desk_card"`
				}{
					Card:     action.card,
					DeskCard: game.GameDeck.card,
				}
				b, _ := json.Marshal(data)
				game.GameDeck.mutex.Unlock()

				game.Broadcast("playCard " + string(b))
			} else {
				game.mutex.Lock()
				seat := strconv.Itoa(game.turn)
				game.mutex.Unlock()
				game.Broadcast("cover by seat " + seat)
			}

		case <-game.finish:
			run = false
			break
		}
	}
}

func (game *Game) RandomPlayer() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(game.Players), func(i, j int) {
		game.Players[i], game.Players[j] = game.Players[j], game.Players[i]
	})

	for k, player := range game.Players {
		if _, ok := player.Card.GetCard()["47"]; ok {
			game.turn = k // who comes first
		}
		player.SetPosition(k)
	}
	sort.Slice(game.Players, func(i, j int) bool {
		return game.Players[i].GetPosition() < game.Players[j].GetPosition()
	})

	fmt.Println("start = " + strconv.Itoa(game.turn))
}

func (game *Game) SendCardStatus() {
	for _, player := range game.Players {
		if !player.isRobot {
			player.User.Write(player.Card)
		}
	}
}

func (game *Game) Broadcast(message any) {
	for _, player := range game.Players {
		if !player.isRobot {
			player.User.Write(message)
		}
	}
}

func (game *Game) CheckTurn(turn int) bool {
	game.mutex.Lock()
	v := game.turn == turn
	fmt.Println(game.turn, turn)
	game.mutex.Unlock()
	return v
}
