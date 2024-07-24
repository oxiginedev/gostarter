package util

import "testing"

func TestIsStringEmpty(t *testing.T) {
	tests := []struct {
		Input string
		Want  bool
	}{
		{
			Input: "something",
			Want:  false,
		}, {
			Input: "",
			Want:  true,
		}, {
			Input: "   ",
			Want:  true,
		},
	}

	for _, test := range tests {
		got := IsStringEmpty(test.Input)
		if got != test.Want {
			t.Errorf("wanted %v but got %v", test.Want, got)
		}
	}
}
