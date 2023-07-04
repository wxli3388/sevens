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
