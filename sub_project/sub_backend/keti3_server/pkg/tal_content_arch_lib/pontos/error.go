package pontos

import "fmt"

type AppError interface {
	error

	String() string
	CodeWithMsg() (code int, msg string)
}

type app struct {
	Cause error
	Code  int
	MSG   string
}

func (a *app) String() string {
	return fmt.Sprintf("caused by %v with code:%d and msg:%s", a.Cause, a.Code, a.MSG)
}

func (a *app) Error() string {
	if a.Cause != nil {
		return a.Cause.Error()
	}

	return ""
}

func (a *app) CodeWithMsg() (code int, msg string) {
	return a.Code, a.MSG
}

// NewAppError generate a wrap error with code and msg
func NewAppError(err error, code int, msg string) AppError {
	return &app{
		Cause: err,
		Code:  code,
		MSG:   msg,
	}
}
