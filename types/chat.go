package types

import (
	"encoding/json"
)

type Chat struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	PoolId    string `json:"poolId"`
	TimeStamp int64  `json:"timeStamp"`
}

type Chats struct {
	ChatRecord []Chat `json:"chatRecord"`
}

func (m Chats) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}
