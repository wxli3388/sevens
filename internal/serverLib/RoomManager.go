package serverLib

import (
	"sync"
)

type JoinRoomInfo struct {
	roomId string
	User   *User
}

type RoomManager struct {
	roomMap map[string]*Room
	userMap map[*User]*Room
	mutex   sync.RWMutex
}

func NewJoinRoomInfo(roomId string, user *User) *JoinRoomInfo {
	return &JoinRoomInfo{
		roomId: roomId,
		User:   user,
	}
}

func NewRoomManager() *RoomManager {
	roomManager := &RoomManager{
		roomMap: map[string]*Room{},
		userMap: map[*User]*Room{},
	}
	return roomManager
}

func (roomManager *RoomManager) createRoom(roomName string) *Room {
	room := NewRoom(roomName)
	roomManager.roomMap[room.roomId] = room
	return room
}

func (roomManager *RoomManager) joinRoom(joinRoomInfo *JoinRoomInfo) bool {
	roomManager.mutex.Lock()
	defer roomManager.mutex.Unlock()
	if joinRoomInfo.roomId == "" {
		for _, room := range roomManager.roomMap {
			if room.JoinRoom(joinRoomInfo.User) {
				roomManager.userMap[joinRoomInfo.User] = room
				return true
			}
		}
		room := roomManager.createRoom("")
		room.JoinRoom(joinRoomInfo.User)
		roomManager.userMap[joinRoomInfo.User] = room
		return true
	}
	if room, ok := roomManager.roomMap[joinRoomInfo.roomId]; ok {
		if room.JoinRoom(joinRoomInfo.User) {
			roomManager.userMap[joinRoomInfo.User] = room
			return true
		}
		return false
	}
	return false
}

func (roomManager *RoomManager) leaveRoom(user *User) bool {
	roomManager.mutex.Lock()
	defer roomManager.mutex.Unlock()
	if room, ok := roomManager.userMap[user]; ok {
		room.LeaveRoom(user)
		user.SetStatus(UserFree)
		delete(roomManager.userMap, user)
		return true
	}
	return false
}

func (roomManager *RoomManager) getUserRoom(user *User) *Room {
	roomManager.mutex.Lock()
	defer roomManager.mutex.Unlock()
	if room, ok := roomManager.userMap[user]; ok {
		return room
	}
	return nil
}

func (roomManager *RoomManager) getAllRoomInfo() []*RoomInfo {
	roomManager.mutex.Lock()
	defer roomManager.mutex.Unlock()
	arr := []*RoomInfo{}
	for _, room := range roomManager.roomMap {
		roomInfo := room.GetRoomInfo()
		arr = append(arr, roomInfo)
	}
	return arr
}
