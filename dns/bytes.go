package dns

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/bits"
	"strings"
)

const (
	ptrCode     byte = 0b11000000
	uint16Bytes int  = 2
)

func write(w io.Writer, values ...uint16) error {
	valueBuf := make([]byte, uint16Bytes)

	for _, value := range values {
		binary.BigEndian.PutUint16(valueBuf, value)

		if _, err := w.Write(valueBuf); err != nil {
			return err
		}
	}

	return nil
}

func readDomain(r *bytes.Reader) (string, error) {
	b := make([]byte, 1)
	sb := strings.Builder{}

	for {
		if _, err := r.Read(b); err != nil {
			return "", err
		}

		length := b[0]

		if length == 0 {
			break
		}

		if length&ptrCode == ptrCode {
			label, err := readCompressedName(r, length)
			if err != nil {
				return "", err
			}

			if sb.Len() > 0 {
				sb.WriteRune('.')
			}

			sb.WriteString(label)

			break
		}

		label := make([]byte, length)
		if _, err := r.Read(label); err != nil {
			return "", err
		}

		if sb.Len() > 0 {
			sb.WriteRune('.')
		}

		sb.Write(label)
	}

	return sb.String(), nil
}

func readCompressedName(r *bytes.Reader, code byte) (string, error) {
	b, err := r.ReadByte()
	if err != nil {
		return "", err
	}

	curOff, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	offset := binary.BigEndian.Uint16([]byte{code & bits.Reverse8(ptrCode), b})

	_, err = r.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return "", err
	}

	result, err := readDomain(r)
	if err != nil {
		return "", err
	}

	_, err = r.Seek(curOff, io.SeekStart)
	if err != nil {
		return "", err
	}

	return result, nil
}
