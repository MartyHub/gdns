package dns

import (
	"bytes"
	"net"
)

const maxResponseSize = 1024

type query struct {
	h *header
	q *question
}

func (q *query) sendRequest(conn net.Conn) error {
	data, err := q.toBytes()
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (q *query) readResponse(conn net.Conn) (*message, error) {
	var err error

	data := make([]byte, maxResponseSize)
	if _, err = conn.Read(data); err != nil {
		return nil, err
	}

	result := new(message)
	err = result.readFrom(bytes.NewReader(data))

	return result, err
}

func (q *query) toBytes() ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := q.h.writeTo(buf); err != nil {
		return nil, err
	}

	if err := q.q.writeTo(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
