package serverLib

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type User struct {
	server     *Server
	connection *websocket.Conn
	send       chan []byte
	status     int
	disConnect chan struct{}
	gameCmd    chan string
}

const (
	UserFree   = 0
	UserInRoom = 1
	UserInGame = 2
)

func NewUser(server *Server, connection *websocket.Conn) *User {
	return &User{
		server:     server,
		connection: connection,
		send:       make(chan []byte, 1024),
		status:     0,
		disConnect: make(chan struct{}),
		gameCmd:    make(chan string),
	}
}

func (user *User) Write(message any) {
	user.send <- []byte(fmt.Sprintf("%v", message))
}

func (user *User) HandleWrite() {
	defer func() {
		user.connection.Close()
	}()

	for {
		select {
		case message, ok := <-user.send:
			user.connection.WriteMessage(websocket.TextMessage, message)
			if !ok {
				return
			}
			fmt.Println(string(message))
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		case <-user.disConnect:
			// fmt.Println("disconnect")
			// close(user.disConnect)
			break
		}
	}
}

func (user *User) HandleRead() {
	defer func() {
		user.connection.Close()
	}()

	for {
		_, message, err := user.connection.ReadMessage()
		if err != nil {
			user.Disconnect()
			break
		}
		msg := string(message)
		if strings.HasPrefix(msg, "game_") {
			if user.status != UserInGame {
				continue
			}
			// room := user.server.roomManager.getUserRoom(user)
			continue
		}
		res := strings.SplitN(msg, " ", 2)
		switch res[0] {
		case "connected":
			roomInfo := user.server.roomManager.getAllRoomInfo()
			cmdRoomInfo := &CmdRoomInfo{RoomInfo: roomInfo}
			user.Write(cmdRoomInfo)
		case "joinRoom":
			if user.status == UserInRoom {
				continue
			}
			roomId := ""
			if len(res) == 2 {
				roomId = res[1]
			}
			res := user.server.roomManager.joinRoom(NewJoinRoomInfo(roomId, user))
			if res {
				user.SetStatus(UserInRoom)
				room := user.server.roomManager.getUserRoom(user)
				if room != nil {
					user.UpdateRoomInfo()
					user.Write(`joinRoom {"roomId": "` + room.roomId + `"}`)
				}
			}
		case "leaveRoom":
			user.server.roomManager.leaveRoom(user)
			user.Write(`leaveRoom {"success": true}`)
			user.UpdateRoomInfo()
		default:

		}

	}
}

func (user *User) UpdateRoomInfo() {
	roomInfo := user.server.roomManager.getAllRoomInfo()
	cmdRoomInfo := &CmdRoomInfo{RoomInfo: roomInfo}
	user.server.Broadcast(cmdRoomInfo)
}

func (user *User) Disconnect() {
	user.server.Unregister <- user
	user.server.roomManager.leaveRoom(user)
	user.disConnect <- struct{}{}
	user.UpdateRoomInfo()
	close(user.send)
	close(user.gameCmd)

}

func (user *User) SetStatus(status int) {
	user.status = status
}