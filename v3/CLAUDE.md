# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PacViz v3 is a terminal user interface (TUI) for browsing and managing Arch Linux packages. It uses the Bubble Tea framework for the TUI and go-alpm for package management.

## Build and Test Commands

```bash
# Build the application
go build ./cmd/pacviz

# Run the application (requires Arch Linux or pacman)
./pacviz

# Run tests
go test ./...

# Run tests for specific packages
go test ./internal/ui/viewport
go test ./internal/command
```

## Architecture

### MVU Pattern (Model-View-Update)

The application follows the Bubble Tea MVU (Model-View-Update) architecture:

- **Model** (`internal/app/model.go`): Central application state including viewport, repository, input modes, and presets
- **Update** (`internal/app/update.go`): Message handlers that update the model based on user input and events
- **View** (`internal/app/view.go`): Renders the UI based on current model state

### Key Components

**Repository Layer** (`internal/repository/`):
- Abstraction over ALPM (Arch Linux Package Management) library
- `Repository` interface defines data access operations
- `AlpmRepository` implements actual ALPM integration
- `MockRepository` for testing

**Domain Layer** (`internal/domain/`):
- `Package`: Complete package metadata including install reason, dependencies, and computed fields
- `Row`: Table row representation derived from packages
- `Preset`: View filters (explicit, dependency, orphans, foreign, all)
- `FilterState`: Text search filtering state

**Viewport** (`internal/ui/viewport/`):
- Manages visible portion of package list with pagination
- Handles selection, sorting, scrolling, and filtering
- Maintains both `AllRows` (full dataset) and `VisibleRows` (filtered subset)
- Separation allows preset filters and text search to work independently

**Column System** (`internal/ui/column/`):
- Column definitions with types, widths, and properties
- Width calculation algorithm: fixed-width columns + auto-expanding Description column
- Columns: Index (#), Repo, Name, Version, Size, Install Date, Groups, Description
- Repo column shows repository name (e.g., "core", "extra") with bright purple color for foreign packages

**Renderer** (`internal/ui/renderer/`):
- Separate rendering components: table, header, status bar, detail panel, scrollbar
- Uses Lipgloss for styling
- Alternating row colors, selection highlighting, sort indicators

**Command System** (`internal/command/`):
- Vim-style command mode (`:quit`, `:preset explicit`, `:goto N`)
- Command parser, registry, and executor
- Commands return `CommandResultMsg` with action results

### Input Modes

The application has three input modes (`internal/app/model.go:14-21`):

1. **Normal Mode**: Navigation, sorting, preset switching
2. **Command Mode**: Vim-style `:` commands
3. **Filter Mode**: Real-time text search with `/`

Mode switching is handled in `internal/app/update.go` with separate handlers for each mode.

### Data Flow

1. **Package Loading**: Repository loads packages → converted to Rows → stored in Viewport
2. **Preset Filtering**: User switches preset (Tab) → applies preset filter → updates VisibleRows
3. **Text Filtering**: User types in filter mode (`/`) → applies text search → updates VisibleRows
4. **Sorting**: User toggles sort (Space) → sorts VisibleRows by selected column
5. **Selection/Scrolling**: User navigates → updates SelectedRow and Offset → re-renders viewport

### Important Patterns

**Preset Application** (`internal/app/model.go:155-166`):
When applying a preset, the code:
1. Resets sort to default (Name, ascending)
2. Clears active text filters
3. Applies the preset filter function

This ensures consistent state when switching between presets.

**Filter Updates** (`internal/ui/viewport/filter.go`):
The viewport maintains separate preset filters (functions) and text filters (terms). Both are applied in sequence to generate VisibleRows.

**Mouse Support** (`internal/app/update.go:226-253`):
Click handling accounts for header row and viewport offset when calculating which package was selected.

## Development Notes

- Package must run on Arch Linux or compatible system with pacman installed
- The repository uses go-alpm which requires ALPM library to be present
- Binary size is ~4.6 MB with debug symbols
- Current phase has 912 packages loading successfully from test system
- The application uses alternative screen mode (`tea.WithAltScreen()`) so it doesn't pollute terminal history

## Testing

Tests exist for:
- Viewport selection logic (`internal/ui/viewport/selection_test.go`)
- Viewport scrolling (`internal/ui/viewport/scroll_test.go`)
- Command execution (`internal/command/execute_test.go`)

When adding new features, maintain test coverage for viewport operations and command handlers.
