package dns

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_question_toBytes(t *testing.T) {
	tests := []struct {
		name string
		q    *question
		want []byte
	}{
		{
			name: "ok",
			q: &question{
				name:    "google.com",
				recType: 5,
				class:   23,
			},
			want: []byte("\x06google\x03com\x00\x00\x05\x00\x17"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := tt.q.writeTo(buf)

			require.NoError(t, err)
			assert.Equal(t, tt.want, buf.Bytes())
		})
	}
}
