package models

type Otp struct {
	ID          uint64 `json:"id,omitempty"`
	WaMessageID string `json:"wa_message_id,omitempty"`
	Phone       string `json:"phone,omitempty"`
	PinCode     string `json:"pin_code,omitempty"`
	Active      bool   `json:"active,omitempty"`
}
