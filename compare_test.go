package geov

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSameMap(t *testing.T) {

	cases := []struct {
		testname string
		first    map[int]bool
		second   map[int]bool
		same     bool
	}{
		{
			testname: "same",
			first:    map[int]bool{1: true, 2: true},
			second:   map[int]bool{1: true, 2: true},
			same:     true,
		},
		{
			testname: "not-same",
			first:    map[int]bool{1: true, 2: true},
			second:   map[int]bool{1: true, 2: true, 3: true},
			same:     false,
		},
		{
			testname: "not-same-2",
			first:    map[int]bool{1: true, 2: true, 3: true},
			second:   map[int]bool{1: true, 2: true},
			same:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.testname, func(t *testing.T) {
			require.Equal(t, tc.same, sameMap(tc.first, tc.second))
		})
	}

}

func TestSubsetMap(t *testing.T) {

	cases := []struct {
		testname string
		set      map[int]bool
		subset   map[int]bool
		isSubset bool
	}{
		{
			testname: "subset-strict",
			set:      map[int]bool{1: true, 2: true},
			subset:   map[int]bool{1: true, 2: true},
			isSubset: true,
		},
		{
			testname: "not-subset",
			set:      map[int]bool{1: true, 2: true},
			subset:   map[int]bool{1: true, 2: true, 3: true},
			isSubset: false,
		},
		{
			testname: "subset",
			set:      map[int]bool{1: true, 2: true, 3: true},
			subset:   map[int]bool{1: true, 2: true},
			isSubset: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.testname, func(t *testing.T) {
			require.Equal(t, tc.isSubset, subsetMap(tc.set, tc.subset))
		})
	}

}
