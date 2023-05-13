package dns

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_header_toBytes(t *testing.T) {
	tests := []struct {
		name string
		h    *header
		want []byte
	}{
		{
			name: "ok",
			h: &header{
				ID:             0x1314,
				Flags:          0,
				NumQuestions:   1,
				NumAnswers:     20,
				NumAuthorities: 300,
				NumAdditionals: 4000,
			},
			want: []byte{'\x13', '\x14', '\x00', '\x00', '\x00', '\x01', '\x00', '\x14', '\x01', '\x2c', '\x0f', '\xa0'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := tt.h.writeTo(buf)

			require.NoError(t, err)
			assert.Equal(t, tt.want, buf.Bytes())
		})
	}
}
