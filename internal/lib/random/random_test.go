package random

import (
	"testing"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "length = 5",
			length: 5,
		},
		{
			name:   "length = 10",
			length: 10,
		},
		{
			name:   "length = 3",
			length: 3,
		},
		{
			name:   "length = 2",
			length: 2,
		},
		{
			name:   "length = 10",
			length: 5,
		},
		{
			name:   "length = 30",
			length: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := make(map[string]struct{})
			for i := 0; i <= 20; i++ {
				res := NewRandomString(tt.length)
				if _, ok := s[res]; ok {
					t.Fatalf("test %v, duplicate string: %v", tt.name, res)
				} else {
					s[res] = struct{}{}
				}
			}
		})
	}
}
