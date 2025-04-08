package commands

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/google/go-cmp/cmp"
)

type testCaseSimple struct {
	name            string
	rawCommand      string
	expectedCommand SimpleParts
}

func TestSimpleCommand(t *testing.T) {
	testCases := []testCaseSimple{
		{
			name:            "empty",
			rawCommand:      "",
			expectedCommand: SimpleParts{Parts: []string{""}},
		},
		{
			name:            "short",
			rawCommand:      "ls",
			expectedCommand: SimpleParts{Parts: []string{"ls"}},
		},
		{
			name:            "3 parts",
			rawCommand:      "gofmt -w ./proto/",
			expectedCommand: SimpleParts{Parts: []string{"gofmt", "-w", "./proto/"}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resultCommand := ParseSimple(testCase.rawCommand)
			if !cmp.Equal(resultCommand.Body, testCase.expectedCommand) {
				t.Errorf("compare failed: %v", deep.Equal(resultCommand.Body, testCase.expectedCommand))
			}
			if resultCommand.String() != testCase.rawCommand {
				t.Errorf("expected %s, got %s", testCase.rawCommand, resultCommand.String())
			}
		})
	}
}
