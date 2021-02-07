package utils

import "testing"

func TestRemoveWhiteSpaces(t *testing.T) {

	type testCase struct {
		before string
		after string
	}

	cases := []testCase{
		{"hello world", "helloworld"},
		{"hello-world", "hello-world"},
		{"-hello-world-", "-hello-world-"},
		{"-hello-world-", "-hello-world-"},
		{"company / repository!1223", "company/repository!1223"},
	}

	for _, c := range cases {
		t.Run(c.before, func(t *testing.T) {
			result := RemoveWhiteSpaces(c.before)
			if c.after != result {
				t.Errorf("RemoveWhiteSpaces('%s') should be equal '%s'. Actual result: '%s'", c.before, c.after, result)
			}
		})
	}
}