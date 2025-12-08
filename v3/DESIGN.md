# PacViz v3 Design Document

## Overview

PacViz is a terminal user interface (TUI) application for managing and exploring Arch Linux packages. It provides an interactive, keyboard-driven alternative to the command-line `pacman` package manager, offering advanced filtering, sorting, searching, and visualization capabilities for both locally installed packages and remote repositories.

## Core Purpose

Provide a fast, efficient, and intuitive TUI for:
- Browsing installed packages with rich metadata
- Discovering orphaned packages and managing dependencies
- Searching remote repositories (official and AUR)
- Installing and removing packages interactively
- Visualizing package relationships and system state

## Design Philosophy

### Guiding Principles

1. **Performance First**: Responsive interaction even with thousands of packages
2. **Progressive Disclosure**: Simple by default, powerful when needed
3. **Keyboard-Centric**: Efficient navigation without touching the mouse
4. **Zero Configuration**: Sensible defaults that work out of the box
5. **Composable Design**: Modular components that can be extended independently
6. **Graceful Degradation**: Handle errors without crashing, inform users clearly

### Design Goals

- **Instant Feedback**: All UI operations should feel instant (<16ms)
- **Minimal Dependencies**: Keep the dependency tree small and auditable
- **Type Safety**: Leverage Go's type system to prevent runtime errors
- **Testability**: Design components to be easily unit tested
- **Maintainability**: Clear separation of concerns, self-documenting code

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Terminal (User)                       │
└────────────────────────┬────────────────────────────────┘
                         │
                    Bubble Tea
                    Event Loop
                         │
           ┌─────────────┼─────────────┐
           │             │             │
        Model ──────── Update ──────── View
           │             │             │
    ┌──────┴──────┐     │      ┌──────┴──────┐
    │             │     │      │             │
Viewport      Pacman  Commands  Renderer   Styles
    │           │       │         │           │
  Rows      Database  Parser   Components  Colors
Columns    Queries   Actions   Layout     Themes
Sorting
Filtering
Selection
```

### Component Layers

1. **Presentation Layer**: View rendering, styles, layout
2. **Application Layer**: Model state, viewport logic, user commands
3. **Domain Layer**: Package data structures, business logic
4. **Data Layer**: Pacman database access, caching

### Key Architectural Patterns

- **Model-View-Update (MVU)**: Unidirectional data flow via Bubble Tea
- **Repository Pattern**: Abstract pacman database access behind interface
- **Command Pattern**: User actions as executable command objects
- **Observer Pattern**: React to data changes (package installations, updates)
- **Virtual Scrolling**: Render only visible rows for performance
- **Lazy Loading**: Load data on-demand, cache intelligently

## Core Features

### 1. Package Browsing

**Purpose**: View and navigate installed packages

**Capabilities**:
- Display packages in a sortable, filterable table
- Show rich metadata: name, version, install date, size, groups, description
- Navigate with keyboard (arrow keys, page up/down, home/end)
- Select packages for actions
- Quick jump to package by name (fuzzy search)

**Data Columns**:
- ID: Row number (computed)
- Name: Package name (sortable, searchable)
- Version: Semantic version (sortable with proper version comparison)
- Size: Installed size in human-readable format (sortable by bytes)
- Install Date: ISO format date (sortable chronologically)
- Groups: Comma-separated package groups (filterable)
- Description: One-line package description (searchable, truncated)

**Views/Presets**:
- **Explicit**: User-installed packages (installation reason = explicit)
- **Dependency**: Auto-installed dependencies
- **Orphans**: Dependencies with no packages depending on them
- **Foreign**: Packages not in any sync database (AUR, manual builds)
- **All**: Every installed package

### 2. Filtering and Searching

**Local Filtering** (instant, in-memory):
- Text filter: Match package names/descriptions (fuzzy or substring)
- Multi-term AND logic: All terms must match
- Column-specific filters: Filter by specific column value
- Regex support: Advanced pattern matching (optional, explicit)
- Case-insensitive by default, case-sensitive toggle

**Remote Search** (query sync databases):
- Search official repositories
- Search AUR (optional integration)
- Display search results in same table interface
- Show installation status (installed/not installed)
- One-key install from search results

**Filter Modes**:
- **Live Filter** (`/`): Real-time filtering as you type
- **Search Mode** (`:search`): Query remote repositories
- **Preset Filters**: Pre-configured views (orphans, explicit, etc.)

### 3. Sorting

**Capabilities**:
- Sort by any column
- Ascending/descending toggle
- Smart sorting by column type:
  - Text: Alphabetical, case-insensitive
  - Version: Semantic version comparison (1.9 < 1.10)
  - Date: Chronological
  - Size: Numerical bytes
- Visual indicator of sort column and direction
- Multi-level sorting (primary + secondary sort key)

**Sort Triggers**:
- Spacebar: Toggle sort on selected column
- Column header click (if mouse support enabled)
- Sort commands: `:sort name asc`, `:sort size desc`

### 4. Package Management

**Installation**:
- Install packages from search results
- Confirm before installation (show size, dependencies)
- Real-time progress display during download/install
- Error handling with actionable messages
- Batch installation (select multiple packages)

**Removal**:
- Remove selected package(s)
- Show reverse dependencies (what will break)
- Cascade removal options (remove with dependencies)
- Confirm before removal with impact summary
- Safe mode (prevent removal of critical packages)

**Updates**:
- Highlight packages with available updates
- Show old version vs new version
- Batch update selection
- Update all (like `pacman -Syu`)

### 5. Package Information

**Detail View**:
- Full package metadata in dedicated pane/modal
- Dependencies (required by this package)
- Reverse dependencies (packages requiring this)
- Optional dependencies with descriptions
- Files installed by package (file list)
- Package URLs (homepage, repository, bug tracker)
- Build date, packager information
- Installation reason, install script presence
- Validation status (signature check)

**Triggers**:
- Enter key on selected package
- `:info` command
- `i` key binding

### 6. Command Mode

**Command Interface**:
- Vim-style command mode (`:` prefix)
- Command autocomplete with suggestions
- Command history (up/down arrow)
- Inline help for commands

**Core Commands**:
- `:q` / `:quit` - Exit application
- `:g <n>` - Go to line number
- `:search <term>` - Search sync databases
- `:install [pkg]` - Install package (current selection if no arg)
- `:remove [pkg]` - Remove package
- `:info [pkg]` - Show package details
- `:sort <col> <asc|desc>` - Sort by column
- `:filter <term>` - Apply filter
- `:preset <name>` - Switch to preset view
- `:help` - Show help screen

### 7. Visual Features

**Display Elements**:
- Header row: Column names with sort indicators (↑↓)
- Data rows: Package information, alternating colors
- Selection highlight: Current row visual indicator
- Status line: Current view, row count, filter status
- Command buffer: Show active command/filter text
- Scrollbar: Visual position indicator in dataset
- Info box: Package details overlay/split pane

**Color Scheme**:
- Dark theme optimized for terminal readability
- Syntax highlighting for package states:
  - Normal: Default text color
  - Orphaned: Warning color (yellow/orange)
  - Foreign: Accent color (purple/magenta)
  - Outdated: Alert color (red)
  - Selected: Highlight background
- Customizable via theme system

**Layout**:
- Responsive to terminal size (handle resize gracefully)
- Adaptive column widths (content-aware, configurable bounds)
- Status bar always visible (bottom)
- Header always visible (top)
- Scrollable content area (middle)
- Optional split pane for info view

### 8. Navigation

**Vertical Navigation**:
- `↑/k` - Move selection up one row
- `↓/j` - Move selection down one row
- `Ctrl+U` - Page up (half screen)
- `Ctrl+D` - Page down (half screen)
- `Home/gg` - Jump to first row
- `End/G` - Jump to last row
- `:g <n>` - Go to specific line number
- `/` + search - Jump to matching package

**Horizontal Navigation** (Column Selection):
- `←/h` - Previous column
- `→/l` - Next column
- Tab - Next column (alternative)
- Shift+Tab - Previous column

**View Navigation**:
- `Tab` - Cycle through presets (Explicit → Orphans → All → ...)
- `:preset <name>` - Jump to specific preset

### 9. Performance Optimization

**Virtual Scrolling**:
- Render only visible rows (viewport height)
- Lazy render buffer (preload rows above/below viewport)
- Efficient diff calculation for re-renders

**Data Caching**:
- Cache pacman database queries
- Invalidate cache on package operations
- Smart cache expiry (time-based + event-based)
- LRU cache for search results

**Async Operations**:
- Background loading for expensive queries
- Loading indicators for async operations
- Cancellable operations (Ctrl+C during search)

**Debouncing**:
- Debounce filter input (wait for typing pause)
- Debounce expensive sorts/filters
- Immediate feedback with optimistic updates

### 10. Error Handling

**Principles**:
- Never panic in production
- Always provide actionable error messages
- Graceful degradation on missing features
- Log errors for debugging

**Error Types**:
- Database errors: Show message, offer retry
- Network errors: Timeout handling, offline mode
- Permission errors: Explain what went wrong, how to fix
- Validation errors: Highlight invalid input, explain why
- System errors: Safe fallback, prevent data corruption

**User Communication**:
- Modal error dialogs for critical errors
- Status line messages for warnings
- Toast notifications for info messages
- Error log accessible via command (`:errors`)

## Data Model

### Core Entities

```go
// Package represents a complete package entry
type Package struct {
    Name           string
    Version        string
    Description    string
    Architecture   string
    URL            string
    Licenses       []string
    Groups         []string
    Dependencies   []string
    OptDepends     map[string]string // name -> description
    Conflicts      []string
    Provides       []string
    Replaces       []string

    // Installation metadata
    Installed      bool
    InstallDate    time.Time
    InstallReason  InstallReason
    InstalledSize  int64
    Packager       string
    BuildDate      time.Time

    // Computed fields
    Required       []string // packages depending on this
    IsOrphan       bool
    IsForeign      bool
    HasUpdate      bool
    NewVersion     string
}

type InstallReason int
const (
    ReasonExplicit InstallReason = iota
    ReasonDependency
)

// Row is the display representation
type Row struct {
    Package  *Package
    Cells    map[ColumnType]string
    Selected bool
    Filtered bool
}

// Column configuration
type Column struct {
    Type        ColumnType
    Name        string
    Width       ColumnWidth
    Sortable    bool
    Searchable  bool
    Visible     bool
}

type ColumnType string
const (
    ColName        ColumnType = "name"
    ColVersion     ColumnType = "version"
    ColSize        ColumnType = "size"
    ColInstallDate ColumnType = "install_date"
    ColGroups      ColumnType = "groups"
    ColDescription ColumnType = "description"
)

type ColumnWidth struct {
    Type WidthType // Fixed, Percent, Auto
    Min  int
    Max  int
    Size int // pixels or percent
}
```

### Application State

```go
type Model struct {
    // Terminal dimensions
    Width  int
    Height int

    // Core components
    Viewport     *Viewport
    PackageRepo  PackageRepository
    CommandMode  *CommandMode
    FilterMode   *FilterMode
    InfoView     *InfoView

    // Application state
    CurrentPreset string
    Error         error
    Loading       bool
}

type Viewport struct {
    // Data
    AllRows      []*Row
    VisibleRows  []*Row

    // Pagination
    Offset       int // First visible row index
    Height       int // Viewport height in rows

    // Selection
    SelectedRow  int
    SelectedCol  int

    // Sorting
    SortColumn   ColumnType
    SortReverse  bool

    // Filtering
    Filter       FilterState

    // Column configuration
    Columns      []*Column
}

type FilterState struct {
    Active bool
    Terms  []string
    Column ColumnType // empty = all columns
    Regex  bool
}
```

## Component Design

### 1. Viewport Component

**Responsibilities**:
- Manage package rows and visibility
- Handle scrolling and pagination
- Apply filters and sorts
- Track selection state
- Calculate visible rows slice

**Key Methods**:
- `Scroll(delta int)`: Move viewport up/down
- `SelectRow(index int)`: Change selection
- `SelectColumn(index int)`: Change column selection
- `ApplyFilter(filter FilterState)`: Filter rows
- `ApplySort(col ColumnType, reverse bool)`: Sort rows
- `GetVisibleRows() []*Row`: Return visible slice
- `GetSelectedPackage() *Package`: Get current package

**Performance Considerations**:
- O(1) scrolling via offset pointer
- Filter and sort results cached
- Lazy evaluation of row visibility
- Reuse row objects (object pooling)

### 2. Package Repository

**Responsibilities**:
- Abstract pacman database access
- Provide query interface for packages
- Cache database results
- Handle database updates/syncs

**Interface**:
```go
type PackageRepository interface {
    GetInstalled() ([]*Package, error)
    GetExplicit() ([]*Package, error)
    GetOrphans() ([]*Package, error)
    GetForeign() ([]*Package, error)
    Search(query string) ([]*Package, error)
    GetPackage(name string) (*Package, error)
    Install(names []string) error
    Remove(names []string, cascade bool) error
    Refresh() error
}
```

**Implementations**:
- `AlpmRepository`: Production implementation using go-alpm
- `MockRepository`: Test implementation for unit tests

**Caching Strategy**:
- Cache installed packages list (invalidate on install/remove)
- Cache search results (TTL: 5 minutes)
- Cache individual package details (LRU, 1000 entries)

### 3. Renderer

**Responsibilities**:
- Convert model state to terminal output
- Handle layout and styling
- Manage responsive column widths
- Render special UI elements (scrollbar, status bar)

**Components**:
- `HeaderRenderer`: Render column headers with sort indicators
- `RowRenderer`: Render package rows with styling
- `StatusRenderer`: Render status/footer line
- `ScrollbarRenderer`: Render scrollbar indicator
- `InfoBoxRenderer`: Render package detail overlay

**Column Width Algorithm**:
1. Calculate required width for each column (max content length)
2. Allocate fixed-width columns first
3. Allocate percentage-based columns
4. Distribute remaining space to auto-sized columns
5. Apply min/max constraints
6. Truncate content exceeding column width (with ellipsis)

### 4. Command Parser

**Responsibilities**:
- Parse user command input
- Validate commands and arguments
- Execute commands on model
- Provide autocomplete suggestions

**Command Structure**:
```go
type Command struct {
    Name        string
    Aliases     []string
    Description string
    Args        []Arg
    Execute     func(m *Model, args []string) tea.Cmd
}

type Arg struct {
    Name     string
    Required bool
    Type     ArgType // String, Int, Enum
    Values   []string // For enum types
}
```

**Command Registry**:
- Commands registered at initialization
- Supports aliases (`:q` = `:quit`)
- Autocomplete from registered commands
- Help text auto-generated from command definitions

### 5. Filter Engine

**Responsibilities**:
- Apply text filters to rows
- Support multiple filter modes
- Optimize filter performance

**Filter Types**:
- Substring match (case-insensitive)
- Exact match (case-sensitive)
- Regex match (advanced)
- Column-specific match

**Filter Application**:
```go
type FilterEngine interface {
    Filter(rows []*Row, filter FilterState) []*Row
    Match(row *Row, filter FilterState) bool
}
```

**Optimization**:
- Compile regex patterns once
- Early exit on first non-match (AND logic)
- Parallel filtering for large datasets (>10k rows)

### 6. Sort Engine

**Responsibilities**:
- Sort rows by column type
- Handle semantic versioning
- Support multi-level sorting

**Sort Functions**:
```go
type SortFunc func(a, b *Row) bool

// Column-specific sort functions
func SortByName(a, b *Row) bool
func SortByVersion(a, b *Row) bool // Semantic version comparison
func SortBySize(a, b *Row) bool
func SortByDate(a, b *Row) bool
```

**Version Sorting**:
- Use semantic version parser (x.y.z-suffix)
- Handle epoch prefixes (1:2.0.0)
- Fall back to string comparison for non-semver

**Multi-Level Sort**:
- Primary sort: Selected column
- Secondary sort: Name (for tie-breaking)
- Stable sort to preserve insertion order

## User Interface Design

### Screen Layout

```
┌─────────────────────────────────────────────────────────────┐
│ ID │ Name        │ Version  │ Size  │ Date       │ Desc...  │ ← Header
├────┼─────────────┼──────────┼───────┼────────────┼──────────┤
│ 1  │ base        │ 3-2      │ 2.3M  │ 2025-01-15 │ Minimal..│
│ 2  │ linux       │ 6.17.7-1 │ 152M  │ 2025-01-20 │ Linux..  │ ← Row
│ 3  │ pacman      │ 7.0.0-1  │ 5.1M  │ 2025-01-18 │ Package..│ ← Selected
│ 4  │ systemd     │ 257.2-1  │ 28M   │ 2025-01-19 │ System.. │
│ ...                                                          │
│                                                              │ ← Scrollable
│                                                              │   Content
│                                                              │
├──────────────────────────────────────────────────────────────┤
│ Explicit (1,234) | Filtered (3) | Sort: Name ↑ | :_        │ ← Status
└──────────────────────────────────────────────────────────────┘
```

### Info View Layout (Split Pane)

```
┌──────────────────────────┬─────────────────────────────────┐
│ ID │ Name    │ Version.. │ Package: pacman                 │
├────┼─────────┼───────────┤ Version: 7.0.0-1                │
│ 1  │ base    │ 3-2       │ Description: Package manager... │
│ 2  │ linux   │ 6.17.7-1  │                                 │
│ 3  │ pacman  │ 7.0.0-1   │ Dependencies:                   │
│ 4  │ systemd │ 257.2-1   │   - glibc                       │
│                           │   - archlinux-keyring           │
│                           │                                 │
│                           │ Required By:                    │
│                           │   - yay                         │
│                           │   - paru                        │
│                           │                                 │
│                           │ Install Date: 2025-01-18        │
│                           │ Install Reason: Explicit        │
│                           │ Installed Size: 5.1 MB          │
│                           │                                 │
├───────────────────────────┴─────────────────────────────────┤
│ Explicit (1,234) | Info View | Press ESC to close          │
└──────────────────────────────────────────────────────────────┘
```

### Key Bindings Reference

**Navigation**:
- `↑/k` - Previous row
- `↓/j` - Next row
- `←/h` - Previous column
- `→/l` - Next column
- `Ctrl+U` - Page up
- `Ctrl+D` - Page down
- `gg/Home` - First row
- `G/End` - Last row

**Actions**:
- `Enter` - Show package info
- `Space` - Toggle sort on column
- `i` - Install package (with confirmation)
- `r` - Remove package (with confirmation)
- `Tab` - Next preset view

**Modes**:
- `/` - Enter filter mode
- `:` - Enter command mode
- `Esc` - Exit mode / close overlay
- `?` - Show help

**Application**:
- `q` / `Ctrl+C` - Quit
- `:q` - Quit (command mode)
- `:help` - Show help screen

## Technical Considerations

### Performance Targets

- **Startup Time**: < 500ms (cold start)
- **Filter Response**: < 16ms (60 FPS)
- **Sort Response**: < 100ms (for 10k packages)
- **Search Response**: < 2s (sync DB search)
- **Render Time**: < 16ms (60 FPS)

### Memory Constraints

- Keep resident memory < 50MB for typical installations (2000 packages)
- Use string interning for repeated strings (package group names, etc.)
- Stream large file lists instead of loading entirely
- Release memory after expensive operations

### Terminal Compatibility

- Support minimum 80x24 terminal size
- Graceful degradation for smaller terminals
- Handle Unicode properly (emoji, special characters)
- Test on common terminals: alacritty, kitty, gnome-terminal, xterm
- Support 256-color and true color modes

### Testing Strategy

**Unit Tests**:
- All core logic components (viewport, filter, sort)
- Package repository interface with mocks
- Command parsing and execution
- Column width calculations

**Integration Tests**:
- Full TUI interactions with simulated input
- Database queries with test fixtures
- Error handling scenarios

**Manual Testing**:
- Terminal compatibility testing
- Performance profiling with large package sets
- Accessibility testing (screen readers, color blindness)

### Build and Distribution

**Build**:
- Single static binary (Go)
- Cross-compile for x86_64 and aarch64
- Strip debug symbols for release builds
- Version embedded at build time

**Distribution**:
- AUR package (PKGBUILD)
- Pre-built binaries on GitHub releases
- Arch Linux official repos (eventual goal)

**Dependencies** (Runtime):
- pacman (required)
- libalpm (required)
- Terminal emulator (required)

## Configuration

### Configuration File

Location: `~/.config/pacviz/config.toml`

```toml
[ui]
theme = "dark"
show_scrollbar = true
mouse_support = false
vim_bindings = true

[columns]
default_visible = ["name", "version", "size", "install_date", "description"]

[columns.name]
width_type = "auto"
min_width = 10
max_width = 40

[columns.version]
width_type = "fixed"
width = 12

[columns.size]
width_type = "fixed"
width = 8

[columns.install_date]
width_type = "fixed"
width = 12

[columns.description]
width_type = "percent"
percent = 40
min_width = 20

[performance]
cache_ttl = 300 # seconds
debounce_filter = 150 # milliseconds
async_search = true

[keybindings]
quit = ["q", "ctrl+c"]
filter = ["/"]
command = [":"]
info = ["enter", "i"]
help = ["?"]

[theme.dark]
accent1 = "#7aa2f7"
accent2 = "#bb9af7"
background = "#1a1b26"
foreground = "#c0caf5"
selected = "#283457"

[pacman]
# Read from /etc/pacman.conf by default
# Override here if needed
# dbpath = "/var/lib/pacman"
```

### CLI Flags

```
pacviz [OPTIONS]

Options:
  -c, --config <path>    Path to config file
  -p, --preset <name>    Start with preset view (explicit, orphans, all)
  -s, --search <term>    Start with search results
  --no-config            Ignore config file, use defaults
  --version              Show version
  --help                 Show help
```

## Future Enhancements

### Planned Features

1. **Dependency Graph Visualization**
   - Show package dependency tree
   - Visual representation of relationships
   - Identify circular dependencies
   - Compute dependency depth

2. **AUR Integration**
   - Search AUR packages
   - Show AUR package info
   - Install from AUR (via helper like yay/paru)
   - Track foreign packages

3. **Package History**
   - Show installation history
   - View package changelog
   - Rollback to previous versions
   - Track package origin (repo vs AUR vs manual)

4. **Advanced Filtering**
   - Save custom filter presets
   - Complex filter expressions (AND/OR/NOT)
   - Filter by package size, date range, group
   - Exclude patterns

5. **Batch Operations**
   - Multi-select packages (visual mode)
   - Batch install/remove
   - Export package lists
   - Import package lists (restore system state)

6. **Themes**
   - Light and dark themes
   - User-defined color schemes
   - Load themes from files

7. **Mouse Support**
   - Click to select rows
   - Click column headers to sort
   - Scroll wheel support
   - Drag to resize columns (if feasible in terminal)

8. **Export Features**
   - Export package list to file
   - Export search results
   - Generate installation script
   - Share package selections

9. **Notifications**
   - Alert on new updates available
   - Background update checks
   - Desktop notifications (optional)

10. **Plugin System**
    - Custom commands
    - Custom views
    - Integration with other tools

### Non-Goals

- **Not a full package manager replacement**: Still use pacman for core operations
- **Not a graphical application**: Terminal-only
- **Not a build system**: Don't compile packages
- **Not a repository manager**: Don't manage pacman mirrors/databases
- **Not cross-platform**: Arch Linux (and derivatives) only

## Success Metrics

**User Experience**:
- Faster than `pacman -Q | grep` for finding packages
- More discoverable than reading pacman man pages
- Lower cognitive load than remembering pacman flags

**Performance**:
- Handle 10,000+ packages smoothly
- Filter response feels instant
- No UI lag during scrolling

**Adoption**:
- GitHub stars > 100
- AUR votes > 50
- Positive feedback from Arch community

## Migration from v1/v2

### What to Keep

**From v1**:
- Bubbletea TUI framework (solid foundation)
- Viewport virtual scrolling concept
- Column width calculation logic (sophisticated, works well)
- Command mode parser (extensible design)
- Filter mode (intuitive UX)
- Preset view system
- Color scheme (professional)

**From v2**:
- Map-based Row structure (cleaner than slice)
- Separated internal package (better organization)
- Global PackageManager singleton (efficient)
- Type-safe column handling with enums
- Viewport sorting/filtering methods

### What to Improve

**Architecture**:
- Repository interface for testability
- Better error handling (no panics)
- Configuration system
- Proper logging

**Features**:
- Implement all planned features (info view, install/remove, etc.)
- Semantic version sorting (fix string-based sorting)
- Adaptive column widths (combine v1 logic + v2 structure)
- Multi-level sorting
- Regex filtering

**Code Quality**:
- Unit tests for core components
- Integration tests for TUI
- Documentation (godoc comments)
- Examples and usage guide

**Performance**:
- Benchmark and optimize hot paths
- Profile memory usage
- Cache expensive operations
- Debounce user input

### Breaking Changes from v1/v2

- Configuration file format (new TOML config)
- CLI flags (standardized)
- Internal API (repository interface)
- Data structures (add fields for new features)

No breaking changes for end users (they just run the binary).

## Conclusion

PacViz v3 represents a ground-up rethinking of the package visualization tool, incorporating lessons learned from v1 and v2. The design focuses on:

- **Performance**: Virtual scrolling, caching, async operations
- **Usability**: Intuitive keyboard navigation, discoverable commands
- **Extensibility**: Plugin system, configuration, themes
- **Reliability**: Error handling, testing, graceful degradation
- **Maintainability**: Clean architecture, documented code, modular design

By adhering to these design principles and implementing the features outlined in this document, PacViz v3 will provide a powerful, efficient, and delightful experience for Arch Linux users managing their packages.
