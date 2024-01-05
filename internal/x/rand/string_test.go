package rand_test

import (
	"testing"

	"github.com/coma/coma/internal/x/rand"
	"github.com/stretchr/testify/assert"
)

func TestRandStr(t *testing.T) {
	testCases := []struct {
		name        string
		expectation int
		actual      func() int
	}{
		{
			name:        "length 64",
			expectation: 64,
			actual: func() int {
				return len(rand.RandStr(64))
			},
		},
		{
			name:        "length 53",
			expectation: 53,
			actual: func() int {
				return len(rand.RandStr(53))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.expectation
			act := tc.actual()

			assert.Equal(t, exp, act)
		})
	}
}
