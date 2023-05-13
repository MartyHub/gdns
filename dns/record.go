package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

//go:generate stringer -type=recordType

type recordType uint16

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.2
const (
	A recordType = iota + 1
	NS
	MD
	MF
	CNAME
	SOA
	MB
	MG
	MR
	NULL
	WKS
	PTR
	HINFO
	MINFO
	MX
	TXT
)

type record struct {
	name    string
	recType recordType
	class   Class
	ttl     uint32
	data    any
}

func (rec *record) dataAsIP() net.IP {
	if rec.recType == A {
		return rec.data.(net.IP) //nolint:forcetypeassert
	}

	return nil
}

func (rec *record) dataAsString() string {
	if rec.recType == CNAME || rec.recType == NS {
		return rec.data.(string) //nolint:forcetypeassert
	}

	if rec.recType == A {
		return rec.data.(net.IP).String() //nolint:forcetypeassert
	}

	return string(rec.data.([]byte)) //nolint:forcetypeassert
}

func (rec *record) readFrom(r *bytes.Reader) error {
	var err error

	if rec.name, err = readDomain(r); err != nil {
		return err
	}

	if err = binary.Read(r, binary.BigEndian, &rec.recType); err != nil {
		return err
	}

	if err = binary.Read(r, binary.BigEndian, &rec.class); err != nil {
		return err
	}

	if err = binary.Read(r, binary.BigEndian, &rec.ttl); err != nil {
		return err
	}

	return rec.readData(r)
}

func (rec *record) readData(r *bytes.Reader) error {
	var (
		err    error
		length uint16
	)

	if err = binary.Read(r, binary.BigEndian, &length); err != nil {
		return err
	}

	if rec.recType == CNAME || rec.recType == NS {
		rec.data, err = readDomain(r)

		return err
	}

	data := make([]byte, length)
	if _, err = r.Read(data); err != nil {
		return err
	}

	if rec.recType == A {
		rec.data = net.IP(data)

		return nil
	}

	rec.data = data

	return nil
}

func (rec *record) String() string {
	return fmt.Sprintf("%s: type=%s class=%s TTL=%v = %s",
		rec.name,
		rec.recType.String(),
		rec.class.String(),
		time.Duration(rec.ttl)*time.Second,
		rec.dataAsString(),
	)
}
