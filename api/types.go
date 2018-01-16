package api

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"`
	status  int    `json:"status,omitempty"`
	//Data    interface{} `json:"data,omitempty"`
}
