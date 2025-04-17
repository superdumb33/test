package entities

type AppErr struct {
	Code    int
	Message string
}

func (e *AppErr) Error() string {
	return e.Message
}

func NewAppErr(code int, message string) *AppErr {
	return &AppErr{Code: code, Message: message}
}