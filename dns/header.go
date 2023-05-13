package dns

import (
	"encoding/binary"
	"io"
)

type header struct {
	ID             uint16
	Flags          uint16
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

func (h *header) readFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, h)
}

func (h *header) writeTo(w io.Writer) error {
	return write(w, h.ID, h.Flags, h.NumQuestions, h.NumAnswers, h.NumAuthorities, h.NumAdditionals)
}
