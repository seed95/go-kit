package derror

import (
	"fmt"
)

type ServerError interface {
	error
	Status() int
	Message() string
	SetDesc(desc string) ServerError
	Description() string
}

type serverError struct {
	code    int
	message string // locale key
	desc    string
}

var _ ServerError = (*serverError)(nil)

func (se *serverError) Error() string {
	if len(se.desc) != 0 {
		return fmt.Sprintf("(status_code: %d) %s, desc: %s", se.code, se.message, se.desc)
	}
	return fmt.Sprintf("(status_code: %d) %s", se.code, se.message)
}

func (se *serverError) Status() int {
	return se.code
}

func (se *serverError) Message() string {
	return se.message
}

func (se *serverError) SetDesc(desc string) ServerError {
	se.desc = desc
	return se
}

func (se *serverError) Description() string {
	return se.desc
}
