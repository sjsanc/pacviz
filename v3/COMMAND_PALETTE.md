# Command Palette Implementation

## Overview

PacViz v3 features a modern command palette that appears when you press `:` in normal mode. The palette shows available commands, filters them as you type, and provides helpful descriptions.

## Features

### Visual Design

When you press `:`, a rounded border box appears above the status bar showing:
- **Command name** (in blue, bold): e.g., `:g`, `:quit`
- **Arguments** (in yellow, italic): e.g., `<line>`, `<col> <asc|desc>`
- **Description** (in gray): Brief explanation of what the command does

Example:
```
╭────────────────────────────────────────────────────────────╮
│ :g <line> - Go to line number                             │
│ :goto <line> - Go to line number                          │
│ :t - Jump to top (first row)                              │
│ :top - Jump to top (first row)                            │
│ :e - Jump to end (last row)                               │
│ :end - Jump to end (last row)                             │
╰────────────────────────────────────────────────────────────╯
:_
```

### Intelligent Filtering

The palette filters commands as you type:
- Press `:` → Shows all commands
- Type `:g` → Shows only commands starting with 'g' (`:g`, `:goto`)
- Type `:t` → Shows only `:t` and `:top`
- Type `:qu` → Shows only `:q` and `:quit`

### Command Aliases

Commands support aliases for convenience:
- `:g` = `:goto`
- `:t` = `:top`
- `:e` = `:end`
- `:q` = `:quit`
- `:?` = `:help`

## Implemented Commands

### Navigation Commands

#### `:g <line>` or `:goto <line>`
Jump to a specific line number (1-based).

**Usage:**
```
:g 100        # Jump to line 100
:goto 50      # Jump to line 50
```

**Behavior:**
- Line numbers are 1-based (user-friendly)
- Automatically bounds-checked (won't go beyond available rows)
- Works with filtered views (jumps within visible rows)

#### `:t` or `:top`
Jump to the first row.

**Usage:**
```
:t           # Jump to top
:top         # Jump to top
```

#### `:e` or `:end`
Jump to the last row.

**Usage:**
```
:e           # Jump to end
:end         # Jump to end
```

### Application Commands

#### `:q` or `:quit`
Exit the application.

**Usage:**
```
:q           # Quit
:quit        # Quit
```

**Note:** You can still use `q` in normal mode or `Ctrl+C` to quit.

## Planned Commands (Not Yet Implemented)

The palette already shows these commands, but they'll be implemented in later phases:

### `:sort <col> <asc|desc>`
Sort by a column in ascending or descending order.

**Example:**
```
:sort name asc
:sort size desc
```

### `:preset <name>`
Switch to a different preset view.

**Example:**
```
:preset explicit
:preset orphans
:preset all
```

### `:filter <term>`
Apply a filter (same as `/` mode).

**Example:**
```
:filter firefox
```

### `:clear`
Clear the current filter.

### `:help` or `:?`
Show the help screen with all keybindings.

## Technical Details

### Command Execution Flow

1. User presses `:` → Enter command mode
2. Command palette renders with all commands
3. User types characters → Palette filters in real-time
4. User presses `Enter` → Command executes
5. Command result processed asynchronously
6. UI updates based on command result

### Command Structure

Commands are defined in `internal/command/palette.go`:

```go
type CommandDef struct {
    Name        string      // Primary command name
    Aliases     []string    // Alternative names
    Args        string      // Argument description
    Description string      // Help text
}
```

### Command Execution

Commands are parsed and executed in `internal/command/execute.go`:

```go
type ExecuteResult struct {
    Quit       bool    // Quit application
    GoToLine   int     // Line to jump to (-1 = no jump)
    ScrollTop  bool    // Scroll to top
    ScrollEnd  bool    // Scroll to end
    Error      string  // Error message
}
```

### Adding New Commands

To add a new command:

1. **Define the command** in `GetAllCommands()` in `palette.go`:
   ```go
   {
       Name:        "mycommand",
       Aliases:     []string{"mc"},
       Args:        "<arg>",
       Description: "Does something cool",
   }
   ```

2. **Implement the handler** in `Execute()` in `execute.go`:
   ```go
   case "mycommand", "mc":
       return executeMyCommand(args)
   ```

3. **Add execution logic**:
   ```go
   func executeMyCommand(args []string) ExecuteResult {
       // Validate args
       // Execute command
       // Return result
   }
   ```

4. **Handle the result** in `handleCommandResult()` in `app/update.go` if needed.

## User Experience

### Discoverability

The command palette improves discoverability:
- **No need to memorize commands**: Press `:` to see all available commands
- **Self-documenting**: Each command shows its arguments and description
- **Fuzzy matching**: Type partial command names to filter
- **Visual feedback**: See what's available before executing

### Consistency

The palette follows modern UI conventions:
- Similar to VS Code's command palette (`Ctrl+Shift+P`)
- Similar to Sublime Text's command palette
- Familiar to users of modern editors

### Performance

- **Instant filtering**: No delay when typing
- **Async execution**: Commands execute without blocking the UI
- **Minimal overhead**: Palette only renders when in command mode

## Keyboard Shortcuts Summary

| Key | Action |
|-----|--------|
| `:` | Open command palette |
| Type | Filter commands in palette |
| `Backspace` | Delete characters from buffer |
| `Enter` | Execute command |
| `Esc` | Cancel and return to normal mode |
| `↑`/`↓` | (Future) Navigate palette items |
| `Tab` | (Future) Autocomplete command |

## Future Enhancements

### Phase 3-4 Enhancements:
- **Arrow key navigation**: Select commands with up/down arrows
- **Tab completion**: Press Tab to autocomplete selected command
- **Command history**: Recall previous commands with up arrow
- **Syntax highlighting**: Highlight command arguments as you type
- **Inline validation**: Show errors before executing (e.g., invalid line number)
- **Command preview**: Show what will happen before executing

### Long-term Enhancements:
- **Custom commands**: User-defined commands via config
- **Command chaining**: Execute multiple commands (`:g 100 | filter linux`)
- **Command scripts**: Save and replay command sequences
- **Fuzzy matching**: Type `:gt` to match `:goto`

## Testing the Command Palette

### Test Scenario 1: Basic Display
1. Launch `./pacviz`
2. Press `:`
3. ✅ Command palette appears with all commands
4. ✅ Status bar shows `:`

### Test Scenario 2: Filtering
1. Press `:`
2. Type `g`
3. ✅ Palette shows only `:g` and `:goto`
4. Type `o`
5. ✅ Palette shows only `:goto`

### Test Scenario 3: Command Execution
1. Press `:`, type `g 100`, press Enter
2. ✅ Jumps to line 100
3. ✅ Palette disappears
4. ✅ Normal mode status bar returns

### Test Scenario 4: Cancel
1. Press `:`, type something
2. Press `Esc`
3. ✅ Returns to normal mode
4. ✅ No command executes

### Test Scenario 5: All Commands
Test each command:
- ✅ `:g 50` → Jumps to line 50
- ✅ `:t` → Jumps to top
- ✅ `:e` → Jumps to end
- ✅ `:q` → Quits application

## Troubleshooting

**Palette doesn't appear:**
- Ensure you're in normal mode (not already in command/filter mode)
- Check terminal size (palette needs minimum height)

**Commands don't execute:**
- Verify command name is correct (check palette)
- Ensure required arguments are provided
- Check for error messages (future feature)

**Filtering doesn't work:**
- Ensure you're typing after the `:` prefix
- Filtering is prefix-based (not fuzzy yet)

## Comparison with v1

**v1:**
- No command palette (just buffer)
- Hard to discover available commands
- Limited commands (`:g`, `:t`, `:e` only)

**v3:**
- Modern command palette UI
- All commands discoverable
- Extensible architecture for adding more commands
- Better user experience
