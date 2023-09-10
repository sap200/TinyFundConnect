package types

type FlowResponseModel struct {
	Result ResultModel `json:"result"`
}

type ResultModel struct {
	NextStep string `json:"next_step"`
}
