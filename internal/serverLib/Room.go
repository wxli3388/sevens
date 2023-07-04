package serverLib

import (
	Utils "sevens/internal/utils"
	"sync"
)

type Room struct {
	roomName      string
	roomId        string
	users         map[*User]bool
	canChangeUser bool
	maxPlayer     int
	mutex         sync.RWMutex
	game          *Game
}

type RoomInfo struct {
	RoomId     string `json:"roomId"`
	RoomName   string `json:"roomName"`
	UsersCount int    `json:"usersCount"`
	CanJoin    bool   `json:"canJoin"`
	MaxPlayer  int    `json:"maxPlayer"`
}

const (
	MaxPlayer = 4
)

func NewRoom(roomName string) *Room {
	if roomName == "" {
		roomName = "Let's play a game"
	}
	return &Room{
		roomName:      roomName,
		roomId:        generateRoomId(),
		users:         make(map[*User]bool),
		maxPlayer:     MaxPlayer,
		canChangeUser: true,
	}
}

func generateRoomId() string {
	return Utils.GenerateRandomString(6)
}

func (room *Room) Broadcast(message string) {
	for user := range room.users {
		user.Write(message)
	}
}

func (room *Room) GetRoomInfo() *RoomInfo {
	return &RoomInfo{
		RoomId:     room.roomId,
		RoomName:   room.roomName,
		UsersCount: len(room.users),
		CanJoin:    room.CanJoin(),
		MaxPlayer:  room.maxPlayer,
	}
}

func (room *Room) CanJoin() bool {
	if len(room.users) < MaxPlayer {
		return true
	}
	return false
}

func (room *Room) StartGame() {
	for user, _ := range room.users {
		user.SetStatus(UserInGame)
	}
	game := NewGame(room.users, room.maxPlayer)
	room.game = game
	// go game.Start()
}

func (room *Room) JoinRoom(user *User) bool {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	if !room.CanJoin() {
		return false
	}
	room.users[user] = true
	if true {
		room.StartGame()
	}
	return true
}

func (room *Room) LeaveRoom(user *User) bool {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	if !room.canChangeUser {
		return false
	}
	if _, ok := room.users[user]; ok {
		delete(room.users, user)
		return true
	}
	return false
}
