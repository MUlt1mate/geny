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
	stringCommand   string
}

func TestProtocCommand(t *testing.T) {
	testCases := []testCaseProtoc{
		{
			name:            "empty",
			rawCommand:      "protoc",
			expectedError:   nil,
			expectedCommand: ProtocBody{},
			stringCommand:   "protoc",
		},
		{
			name:          "variation 1",
			rawCommand:    "protoc file.proto",
			expectedError: nil,
			expectedCommand: ProtocBody{
				Files: []string{"file.proto"},
			},
			stringCommand: "protoc file.proto",
		},
		{
			name:          "all variables",
			rawCommand:    "protoc -I=. -I=./vendor --httpgo_out=. --httpgo_opt=paths=source_relative,marshaller=easyjson proto/example.proto proto/example2.proto",
			expectedError: nil,
			expectedCommand: ProtocBody{
				Imports: []string{".", "./vendor"},
				Plugins: []ProtocPlugin{
					{
						Name: "httpgo",
						Parameters: []ProtocPluginKV{
							{Name: "paths", Value: "source_relative"},
							{Name: "marshaller", Value: "easyjson"},
						},
					},
				},
				Files: []string{"proto/example.proto", "proto/example2.proto"},
			},
			stringCommand: "protoc -I=. -I=./vendor --httpgo_out=paths=source_relative,marshaller=easyjson:. proto/example.proto proto/example2.proto",
		},
		{
			name:          "empty spaces",
			rawCommand:    "  protoc -I=. -I=./vendor --httpgo_out=.		   --httpgo_opt=paths=source_relative,marshaller=easyjson proto/example.proto  	proto/example2.proto  	 ",
			expectedError: nil,
			expectedCommand: ProtocBody{
				Imports: []string{".", "./vendor"},
				Plugins: []ProtocPlugin{
					{
						Name: "httpgo",
						Parameters: []ProtocPluginKV{
							{Name: "paths", Value: "source_relative"},
							{Name: "marshaller", Value: "easyjson"},
						},
					},
				},
				Files: []string{"proto/example.proto", "proto/example2.proto"},
			},
			stringCommand: "protoc -I=. -I=./vendor --httpgo_out=paths=source_relative,marshaller=easyjson:. proto/example.proto proto/example2.proto",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resultCommand, resultError := ParseProtoc(testCase.rawCommand)
			if !cmp.Equal(resultError, testCase.expectedError) {
				t.Errorf("expected error %v got: %v", testCase.expectedError, resultError)
			}
			if !cmp.Equal(resultCommand.Body, testCase.expectedCommand) {
				t.Errorf("compare failed: %v", deep.Equal(resultCommand.Body, testCase.expectedCommand))
			}
			if resultCommand.String() != testCase.stringCommand {
				t.Errorf("expected string command:\n%s\ngot:\n%s", testCase.stringCommand, resultCommand.String())
			}
		})
	}
}
