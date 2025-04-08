package commands

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/google/go-cmp/cmp"
)

type testCaseProtoc struct {
	name            string
	rawCommand      string
	expectedError   error
	expectedCommand ProtocBody
}

func TestProtocCommand(t *testing.T) {
	testCases := []testCaseProtoc{
		{
			name:            "empty",
			rawCommand:      "protoc",
			expectedError:   nil,
			expectedCommand: ProtocBody{},
		},
		{
			name:          "variation 1",
			rawCommand:    "protoc file.proto",
			expectedError: nil,
			expectedCommand: ProtocBody{
				Files: []string{"file.proto"},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resultCommand, resultError := ParseProtoc(testCase.rawCommand)
			if !cmp.Equal(resultError, testCase.expectedError) {
				t.Errorf("compare failed: %v", deep.Equal(resultError, testCase.expectedError))
			}
			if !cmp.Equal(resultCommand.Body, testCase.expectedCommand) {
				t.Errorf("compare failed: %v", deep.Equal(resultCommand.Body, testCase.expectedCommand))
			}
			if resultCommand.String() != testCase.rawCommand {
				t.Errorf("expected %s, got %s", testCase.rawCommand, resultCommand.String())
			}
		})
	}
}
