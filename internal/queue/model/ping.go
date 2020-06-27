package model

// PingPayload for ping queue messages
type PingPayload struct {
	Msg string `json:"message"`
}
