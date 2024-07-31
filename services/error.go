package services

type ServiceError struct {
	ErrCode  int    `json:"status_code"`
	ErrMsg   string `json:"message"`
	Internal error  `json:"-"`
}

func (se *ServiceError) Error() string {
	return se.ErrMsg
}

func NewServiceError(code int, msg string, err error) *ServiceError {
	return &ServiceError{
		ErrCode:  code,
		ErrMsg:   msg,
		Internal: err,
	}
}
