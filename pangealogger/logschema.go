package pangealogger

// Custom schema implementation

type CustomSchemaEvent struct {
	Message   string `json:"message"`
	Event     string `json:"event,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
	PoolId    string `json:"pool_id,omitempty"`

	// TenantID field
	TenantID string `json:"tenant_id,omitempty"`
}

func (e *CustomSchemaEvent) GetTenantID() string {
	return e.TenantID
}

func (e *CustomSchemaEvent) SetTenantID(tid string) {
	e.TenantID = tid
}

func New(message, event, ipAddress, email, poolId string) *CustomSchemaEvent {
	return &CustomSchemaEvent{
		Message:   message,
		Event:     event,
		IPAddress: ipAddress,
		UserEmail: email,
		PoolId:    poolId,
	}
}
