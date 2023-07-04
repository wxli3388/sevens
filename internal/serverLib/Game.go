package serverLib

import (
	"math/rand"
	"time"
)

type Game struct {
	Players []*Player

	Action    chan string
	finish    chan struct{}
	turn      int
	maxPlayer int
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
		Action:    make(chan string),
		finish:    make(chan struct{}),
		maxPlayer: maxPlayer,
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

	// IntervalTime := 5 * time.Second        // 觸發間隔時間
	// ticker := time.NewTicker(IntervalTime) // 設定 5 秒觸發一次
	game.Players[game.turn].RoundStart()
	for {
		select {
		case <-game.Action:
		// 	fmt.Println(act)
		// case c := <-ticker.C:
		// 	game.turn = (game.turn + 1) % game.maxPlayer
		// fmt.Println("now: ", c)
		case <-game.finish:
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
