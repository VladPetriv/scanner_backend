package convert_test

import (
	"testing"

	"github.com/VladPetriv/scanner_backend/pkg/convert"
	"github.com/stretchr/testify/assert"
)

func Test_PageToOffset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "PageToOffset converts 0",
			input: 0,
			want:  0,
		},
		{
			name:  "PageToOffset converts 1",
			input: 1,
			want:  0,
		},
		{
			name:  "PageToOffset converts 10",
			input: 10,
			want:  90,
		},
		{
			name:  "PageToOffset converts 2",
			input: 2,
			want:  10,
		},
		{
			name:  "PageToOffset converts 30",
			input: 30,
			want:  290,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := convert.PageToOffset(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
