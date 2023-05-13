package dns

import (
	"bytes"
	"encoding/binary"
	"strings"
)

const maxLabelLength = 63

type question struct {
	name    string
	recType recordType
	class   Class
}

func (q *question) readFrom(r *bytes.Reader) error {
	var err error

	if q.name, err = readDomain(r); err != nil {
		return err
	}

	if err = binary.Read(r, binary.BigEndian, &q.recType); err != nil {
		return err
	}

	//nolint:revive
	if err = binary.Read(r, binary.BigEndian, &q.class); err != nil {
		return err
	}

	return nil
}

func (q *question) writeTo(buf *bytes.Buffer) error {
	if err := q.writeName(buf); err != nil {
		return err
	}

	return write(buf, uint16(q.recType), uint16(q.class))
}

func (q *question) writeName(buf *bytes.Buffer) error {
	for _, label := range strings.Split(q.name, ".") {
		l := len(label)

		if l > maxLabelLength {
			return NewError("label too long: %s", label)
		}

		buf.WriteByte(byte(l))
		buf.WriteString(label)
	}

	buf.WriteByte(0)

	return nil
}
