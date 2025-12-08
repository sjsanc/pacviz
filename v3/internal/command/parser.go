package command

import "strings"

// Parse parses a command string and returns the command name and arguments.
func Parse(input string) (string, []string) {
	// Remove leading `:` if present
	input = strings.TrimPrefix(input, ":")
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return "", nil
	}

	return parts[0], parts[1:]
}

// Match finds a command by name or alias.
func Match(name string, registry *Registry) *Command {
	// TODO: Implement command matching logic
	// TODO: Support aliases
	return registry.Get(name)
}
