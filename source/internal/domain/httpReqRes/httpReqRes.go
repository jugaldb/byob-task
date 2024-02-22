package httpReqRes

type HttpResponse struct {
	Data    any    `json:"data"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	RESPONSE_SUCCESS string = "SUCCESS"
	RESPONSE_FAILURE string = "FAILURE"
)
