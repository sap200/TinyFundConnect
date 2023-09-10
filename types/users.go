package types

import "encoding/json"

type User struct {
	Email         string   `json:"email"`
	PasswordHash  string   `json:"passwordHash"`
	Password      string   `json:"password"`
	CustomerId    string   `json:"customerId"`
	EmailVerified bool     `json:"emailVerified"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	Ip            string   `json:"ip"`
	FlowId        string   `json:"flowId"`
	Code          string   `json:"code"`
	Cancel        bool     `json:"cancel"`
	PoolIds       []string `json:"poolIds"`
	Balance       float64  `json:"balance"`
}

func (u User) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}
