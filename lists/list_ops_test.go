package lists

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListOpsValidateIndex(t *testing.T) {
	type testCase struct {
		index    int
		size     int
		expected string // Expected error message or empty string for no error
	}

	// Define test cases
	cases := []testCase{
		// Valid indices
		{index: 0, size: 5, expected: ""},
		{index: 4, size: 5, expected: ""},

		// Invalid indices (negative index)
		{index: -1, size: 5, expected: "index out of bounds, must be between 0 and 4, but -1 was provided"},

		// Invalid indices (greater than or equal to size)
		{index: 5, size: 5, expected: "index out of bounds, must be between 0 and 4, but 5 was provided"},
		{index: 10, size: 5, expected: "index out of bounds, must be between 0 and 4, but 10 was provided"},
	}

	// Iterate through each test case
	for _, tc := range cases {
		t.Run(fmt.Sprintf("index=%d,size=%d", tc.index, tc.size), func(t *testing.T) {
			ops := listOps[int]{}
			err := ops.validateIndex(tc.index, tc.size)

			// If expected is an empty string, we assert there should be no error
			if tc.expected == "" {
				assert.NoError(t, err, "unexpected error for index=%d and size=%d", tc.index, tc.size)
			} else {
				// Assert the error is not nil and assert the error message
				assert.Error(t, err, "expected an error for index=%d and size=%d", tc.index, tc.size)
				assert.EqualError(t, err, tc.expected, "unexpected error message")
			}
		})
	}
}
