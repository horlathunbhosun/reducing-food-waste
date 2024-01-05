package response

type JsonResponse struct {
	Error        bool        `json:"error"`
	ErrorMessage interface{} `json:"error_message,omitempty"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
	Status       bool        `json:"status"`
}
