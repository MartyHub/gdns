package dns

import "fmt"

type Error struct {
	Msg string
}

func NewError(format string, a ...any) Error {
	return Error{Msg: "gdns: " + fmt.Sprintf(format, a...)}
}

func (e Error) Error() string {
	return e.Msg
}
