package util

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func BuildErrorResponse(msg string, errors interface{}) Response {
	return newResponse(false, msg, nil, errors)
}

func BuildSuccessResponse(msg string, data interface{}) Response {
	return newResponse(true, msg, data, nil)
}

func newResponse(status bool, msg string, data, errors interface{}) Response {
	return Response{
		Status:  status,
		Message: msg,
		Data:    data,
		Errors:  errors,
	}
}
