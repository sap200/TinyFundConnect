package types

import "encoding/json"

type Pool struct {
	PoolName     string   `json:"poolName"`
	PoolId       string   `json:"poolId"`
	PoolBalance  float64  `json:"poolBalance"`
	CreatorEmail string   `json:"creatorEmail"`
	MembersList  []Member `json:"membersList"`
	OneMember    Member   `json:"oneMember"`
}

func (p Pool) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}

type Member struct {
	Email     string  `json:"email"`
	Balance   float64 `json:"balance"`
	HasJoined bool    `json:"hasJoined"`
}

func (m Member) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}

type UserWithPoolAndStatus struct {
	IsActive    bool `json:"isActive"`
	CurrentPool Pool `json:"currentPool"`
}

type PoolListWithUserMail struct {
	UserEmail string                  `json:"userEmail"`
	PoolList  []UserWithPoolAndStatus `json:"poolList"`
}

func (m PoolListWithUserMail) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}
