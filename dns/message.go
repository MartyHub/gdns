package dns

import (
	"bytes"
	"net"
	"strings"
)

type message struct {
	header      header
	questions   []question
	answers     []record
	authorities []record
	additionals []record
}

func (m *message) answer() (*record, error) {
	if len(m.answers) != 1 {
		return nil, NewError("expected 1 answer, got %d", len(m.answers))
	}

	return &m.answers[0], nil
}

func (m *message) firstIP() (net.IP, error) {
	for _, auth := range m.authorities {
		if auth.recType == A {
			return auth.dataAsIP(), nil
		}
	}

	for _, auth := range m.additionals {
		if auth.recType == A {
			return auth.dataAsIP(), nil
		}
	}

	return nil, NewError("no IP in %s", m.String())
}

func (m *message) readFrom(r *bytes.Reader) error {
	var err error

	if err = m.readHeader(r); err != nil {
		return err
	}

	if err = m.readQuestions(r); err != nil {
		return err
	}

	if m.answers, err = readRecords(r, m.header.NumQuestions); err != nil {
		return err
	}

	if m.authorities, err = readRecords(r, m.header.NumAuthorities); err != nil {
		return err
	}

	if m.additionals, err = readRecords(r, m.header.NumAdditionals); err != nil {
		return err
	}

	return nil
}

func (m *message) readHeader(r *bytes.Reader) error {
	return (&m.header).readFrom(r)
}

func (m *message) readQuestions(r *bytes.Reader) error {
	for i := uint16(0); i < m.header.NumQuestions; i++ {
		q := new(question)
		if err := q.readFrom(r); err != nil {
			return err
		}

		m.questions = append(m.questions, *q)
	}

	return nil
}

func (m *message) String() string {
	sb := strings.Builder{}

	sb.WriteString("Message:")

	for _, a := range m.answers {
		sb.WriteString("\n\tanswer: ")
		sb.WriteString(a.String())
	}

	for _, a := range m.authorities {
		sb.WriteString("\n\tauthority: ")
		sb.WriteString(a.String())
	}

	for _, a := range m.additionals {
		sb.WriteString("\n\tadditional: ")
		sb.WriteString(a.String())
	}

	return sb.String()
}

func readRecords(r *bytes.Reader, num uint16) ([]record, error) {
	var result []record

	for i := uint16(0); i < num; i++ {
		rec := new(record)
		if err := rec.readFrom(r); err != nil {
			return result, err
		}

		result = append(result, *rec)
	}

	return result, nil
}
