package dns

//go:generate stringer -type=Class

type Class uint16

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.4
const (
	IN Class = iota + 1
	CS
	CH
	HS
)
