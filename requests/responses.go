package requests

type WhatsappSendMessageErrorResponse struct {
	ErrorResponse WhatsappSendMessageErrorContent `json:"error"`
}

type WhatsappSendMessageErrorContent struct {
	Message      string                       `json:"message"`
	Type         string                       `json:"type"`
	Code         string                       `json:"code"`
	ErrorData    WhatsappSendMessageErrorData `json:"error_data"`
	ErrorSubcode string                       `json:"error_subcode"`
	FbTraceId    string                       `json:"fbtrace_id"`
}

type WhatsappSendMessageErrorData struct {
	MessagingProduct string `json:"messaging_product"`
	Details          string `json:"details"`
}

type WhatsappSendMessageResponse struct {
	MessagingProduct string            `json:"messaging_product",omitemtpy`
	Contacts         []WhatsappContact `json:"contacts",omitempty`
	Messages         []WhatsappMessage `json:"messages",omitempty`
}

type WhatsappContact struct {
	Input string `json:"input",omitempty`
	WaId  string `json:"wa_id",omitempty`
}

type WhatsappMessage struct {
	ID string `json:"id",omitempty`
}
