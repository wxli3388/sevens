package serverLib

import (
	"encoding/json"
)

type CmdRoomInfo struct {
	RoomInfo []*RoomInfo `json:"roomInfo"`
}

func (cri *CmdRoomInfo) String() string {
	s, err := json.Marshal(cri)
	if err != nil {
		return "roomInfo {}"
	}
	return "roomInfo " + string(s)
}

type CmdOutJoinRoom struct {
	RoomId string `json:"roomId"`
}

func (c *CmdOutJoinRoom) String() string {
	s, err := json.Marshal(c)
	if err != nil {
		return "joinRoom {}"
	}
	return "joinRoom " + string(s)
}

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
