package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_numRange(t *testing.T) {
	tests := map[string]struct {
		min  int
		max  int
		want []int
	}{
		"single element": {
			min:  3,
			max:  3,
			want: []int{3},
		},
		"multiple elements": {
			min:  0,
			max:  7,
			want: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := numRange(tc.min, tc.max)

			assert.Equal(t, tc.want, got)
		})
	}
}
