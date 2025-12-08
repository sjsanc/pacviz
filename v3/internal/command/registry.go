package command

// Registry stores all available commands.
type Registry struct {
	commands map[string]*Command
}

// NewRegistry creates a new command registry.
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]*Command),
	}
}

// Register adds a command to the registry.
func (r *Registry) Register(cmd *Command) {
	r.commands[cmd.Name] = cmd
	// TODO: Register aliases
}

// Get retrieves a command by name.
func (r *Registry) Get(name string) *Command {
	return r.commands[name]
}

// All returns all registered commands.
func (r *Registry) All() []*Command {
	cmds := make([]*Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// Autocomplete returns command suggestions for a partial input.
func (r *Registry) Autocomplete(partial string) []string {
	// TODO: Implement autocomplete logic
	return nil
}
