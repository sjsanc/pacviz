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

func TestExecute_Preset(t *testing.T) {
	tests := []struct {
		name           string
		commandStr     string
		expectedPreset string
		expectedError  string
	}{
		{
			name:           "preset with p alias - explicit",
			commandStr:     "p explicit",
			expectedPreset: "explicit",
			expectedError:  "",
		},
		{
			name:           "preset with full name - dependency",
			commandStr:     "preset dependency",
			expectedPreset: "dependency",
			expectedError:  "",
		},
		{
			name:           "preset orphans",
			commandStr:     "p orphans",
			expectedPreset: "orphans",
			expectedError:  "",
		},
		{
			name:           "preset foreign",
			commandStr:     "p foreign",
			expectedPreset: "foreign",
			expectedError:  "",
		},
		{
			name:           "preset all",
			commandStr:     "p all",
			expectedPreset: "all",
			expectedError:  "",
		},
		{
			name:           "preset aur",
			commandStr:     "p aur",
			expectedPreset: "aur",
			expectedError:  "",
		},
		{
			name:           "preset updatable",
			commandStr:     "p updatable",
			expectedPreset: "updatable",
			expectedError:  "",
		},
		{
			name:           "preset without args",
			commandStr:     "p",
			expectedPreset: "",
			expectedError:  "Usage: :p <preset> (explicit, dependency, orphans, foreign, aur, updatable, all)",
		},
		{
			name:           "preset with invalid name",
			commandStr:     "p invalid",
			expectedPreset: "",
			expectedError:  "Invalid preset: invalid (valid: explicit, dependency, orphans, foreign, aur, updatable, all)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Execute(tt.commandStr)

			if result.PresetChange != tt.expectedPreset {
				t.Errorf("PresetChange = %q, want %q", result.PresetChange, tt.expectedPreset)
			}

			if result.Error != tt.expectedError {
				t.Errorf("Error = %q, want %q", result.Error, tt.expectedError)
			}
		})
	}
}
