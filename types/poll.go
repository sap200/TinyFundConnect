package types

import (
	"encoding/json"
)

type Poll struct {
	PoolId           string   `json:"poolId"`
	PollId           string   `json:"pollId"`
	Message          string   `json:"message"`
	IsActive         bool     `json:"isActive"`
	PollCreatorEmail string   `json:"pollCreatorEmail"`
	Options          []string `json:"options"`
	Votes            []Vote   `json:"votes"`
}

type Vote struct {
	PoolId     string `json:"poolId"`
	PollId     string `json:"pollId"`
	Choice     int    `json:"choice"`
	VoterEmail string `json:"voterEmail"`
}

func (v Vote) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}

func (p Poll) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}

func (p Poll) GetResults() (map[int]int, error) {

	m := make(map[int]int)
	for i := 0; i < len(p.Options); i++ {
		m[i] = 0
	}

	if p.Votes == nil {
		return m, nil
	}

	for i := 0; i < len(p.Votes); i++ {
		m[p.Votes[i].Choice] += 1
	}

	// returns choiceIndex, voteCount, OptionValue
	return m, nil
}
