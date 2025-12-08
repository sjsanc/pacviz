package command

import "testing"

func TestExecute_GoTo(t *testing.T) {
	tests := []struct {
		name           string
		commandStr     string
		expectedLine   int
		expectedError  string
	}{
		{
			name:         "goto with g alias",
			commandStr:   "g 100",
			expectedLine: 99, // 0-indexed
			expectedError: "",
		},
		{
			name:         "goto with full name",
			commandStr:   "goto 50",
			expectedLine: 49,
			expectedError: "",
		},
		{
			name:         "goto line 1",
			commandStr:   "g 1",
			expectedLine: 0,
			expectedError: "",
		},
		{
			name:         "goto without args",
			commandStr:   "g",
			expectedLine: -1,
			expectedError: "Usage: :g <line>",
		},
		{
			name:         "goto with invalid number",
			commandStr:   "g abc",
			expectedLine: -1,
			expectedError: "Invalid line number: abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Execute(tt.commandStr)

			if result.GoToLine != tt.expectedLine {
				t.Errorf("GoToLine = %d, want %d", result.GoToLine, tt.expectedLine)
			}

			if result.Error != tt.expectedError {
				t.Errorf("Error = %q, want %q", result.Error, tt.expectedError)
			}
		})
	}
}
