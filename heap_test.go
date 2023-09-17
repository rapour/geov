package geov

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Number int

func (n Number) Value() int {
	return int(n)
}

func TestHeap(t *testing.T) {

	cases := []struct {
		testname    string
		input       []Element[int]
		sorted      []int
		expectedMin Number
	}{
		{
			testname: "first case",
			input:    []Element[int]{Number(1), Number(3), Number(5), Number(0)},
			sorted:   []int{0, 1, 3, 5},

			expectedMin: 0,
		},
		{
			testname:    "second case",
			input:       []Element[int]{Number(11), Number(3), Number(5), Number(10)},
			sorted:      []int{3, 5, 10, 11},
			expectedMin: 3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.testname, func(t *testing.T) {

			h := NewHeap(tc.input)
			h.BuildMinHeap()
			require.Equal(t, tc.expectedMin, h.Min())

		})
	}

	for _, tc := range cases {
		t.Run(tc.testname, func(t *testing.T) {

			h := NewHeap(tc.input)
			h.BuildMinHeap()

			var sorted []int
			for {
				if h.GetSize() == 0 {
					break
				}
				sorted = append(sorted, h.ExtractMin().Value())
			}
			require.Equal(t, tc.sorted, sorted)

		})
	}

}
