# PacViz v3 Interactions Guide

This document outlines all user interactions in PacViz, including implemented features and planned enhancements.

## Status Legend
- ✅ **Implemented**: Feature is complete and working
- 🚧 **In Progress**: Currently being implemented
- 📋 **Planned**: Designed but not yet implemented

---

## Navigation

### Vertical Navigation
- ✅ `↑` / `k` - Move selection up one row
- ✅ `↓` / `j` - Move selection down one row
- ✅ `Ctrl+U` - Page up (half screen)
- ✅ `Ctrl+D` - Page down (half screen)
- ✅ `Home` / `g` - Jump to first row
- ✅ `G` / `End` - Jump to last row
- 📋 `gg` - Jump to first row (vim-style double tap)

### Horizontal Navigation (Column Selection)
- ✅ `←` / `h` - Previous column
- ✅ `→` / `l` - Next column
- 📋 `Tab` - Next column (alternative)
- 📋 `Shift+Tab` - Previous column

### View Navigation
- 📋 `Tab` - Cycle through presets (Explicit → Orphans → All → ...)
- 📋 `:preset <name>` - Jump to specific preset

---

## Input Modes

### Normal Mode (Default)
The default mode for navigating and viewing packages.

**Available Keys**:
- All navigation keys (see above)
- `q` / `Ctrl+C` - Quit application
- `:` - Enter command mode
- `/` - Enter filter mode
- `Space` - Toggle sort on selected column
- `Enter` / `i` - Show package info
- `?` - Show help screen

### Command Mode
🚧 **Status**: Phase 1 - In Progress

Activated by pressing `:` in normal mode. Provides a vim-style command interface.

**Behavior**:
- Command buffer appears in status line (bottom of screen)
- Type commands and press `Enter` to execute
- Press `Esc` to cancel and return to normal mode
- Press `Backspace` to delete characters

**Available Commands**:

#### Navigation Commands
- 📋 `:g <n>` - Go to line number `n`
- 📋 `:t` or `:top` - Jump to top (first row)
- 📋 `:e` or `:end` - Jump to end (last row)

#### View Commands
- 📋 `:preset <name>` - Switch to preset view
  - Options: `explicit`, `dependency`, `orphans`, `foreign`, `all`
- 📋 `:sort <column> <direction>` - Sort by column
  - Columns: `name`, `version`, `size`, `date`, `groups`, `description`
  - Direction: `asc` or `desc`

#### Filter Commands
- 📋 `:filter <term>` - Apply filter (same as `/` mode)
- 📋 `:clear` - Clear current filter

#### Package Management
- 📋 `:install [pkg]` - Install package (current selection if no arg)
- 📋 `:remove [pkg]` - Remove package (current selection if no arg)
- 📋 `:info [pkg]` - Show package details (current selection if no arg)

#### Search Commands
- 📋 `:search <term>` - Search sync databases (official repos)
- 📋 `:aur <term>` - Search AUR packages

#### Application Commands
- 📋 `:q` or `:quit` - Exit application
- 📋 `:help` - Show help screen
- 📋 `:refresh` - Refresh package database

**Command History**:
- 📋 `↑` - Previous command in history
- 📋 `↓` - Next command in history

**Autocomplete**:
- 📋 `Tab` - Autocomplete command name
- 📋 `Tab` - Cycle through command suggestions

### Filter Mode
🚧 **Status**: Phase 1 - In Progress

Activated by pressing `/` in normal mode. Provides live filtering of the package list.

**Behavior**:
- Filter buffer appears in status line with `/` prefix
- Filters update in real-time as you type
- Matches package names and descriptions (case-insensitive substring match)
- Press `Enter` to accept filter and hide buffer (keeps filter active)
- Press `Esc` to cancel and reset filter (returns to unfiltered view)
- Press `Backspace` to delete characters

**Filter Syntax**:
- Simple text: `firefox` - matches packages containing "firefox"
- Multiple terms (AND): 📋 `python gtk` - matches packages containing both terms
- Column-specific: 📋 `name:firefox` - matches only package names
- Regex mode: 📋 `/^firefox/` - regex pattern matching

**Filter Indicators**:
- Status line shows: `Filter: "<term>" | Showing X of Y packages`

---

## Sorting

### Interactive Sorting
- 📋 `Space` - Toggle sort on currently selected column
  - First press: Sort ascending
  - Second press: Sort descending
  - Third press: Remove sort (restore default order)

### Visual Indicators
- 📋 Column header shows sort direction: `Name ↑` (ascending) or `Name ↓` (descending)
- 📋 Sorted column is highlighted in header

### Sort Behavior
- **Smart sorting by type**:
  - Text columns: Alphabetical, case-insensitive
  - Version column: Semantic version comparison (1.9 < 1.10 < 2.0)
  - Size column: Numerical byte comparison
  - Date column: Chronological
- **Secondary sort**: Always sorts by name as tiebreaker
- **Stable sort**: Preserves order for equal elements

---

## Package Information

### Info View
📋 **Status**: Phase 3 - Planned

Activated by pressing `Enter` or `i` on a selected package.

**Display Modes**:

#### Split Pane Mode (Default)
- Screen splits vertically (60/40)
- Left: Package list (navigable)
- Right: Package details (scrollable)

#### Modal/Overlay Mode (Alternative)
- Full-screen overlay with package details
- Background dimmed

**Info View Content**:
- Package name and version
- Description (full, not truncated)
- Architecture
- URL, licenses
- Groups
- Dependencies (required by this package)
- Reverse dependencies (packages requiring this)
- Optional dependencies with descriptions
- Conflicts, provides, replaces
- Install date, install reason, installed size
- Packager, build date
- Validation status
- File list (expandable section)

**Info View Controls**:
- `Esc` - Close info view and return to list
- `↑` / `↓` - Scroll info content
- `q` - Close info view
- `Tab` - Switch between list and info pane
- `Enter` - Jump to dependency (if focused on dependency name)

---

## Preset Views

📋 **Status**: Phase 3 - Planned

Preset views filter packages based on common use cases.

### Available Presets
1. **Explicit** (Default)
   - Packages explicitly installed by user
   - Installation reason = explicit
   - Most common view for managing installed packages

2. **Dependency**
   - Packages installed as dependencies
   - Installation reason = dependency
   - Useful for seeing automatically installed packages

3. **Orphans**
   - Dependencies with no packages depending on them
   - Can often be safely removed
   - Helps clean up unused packages

4. **Foreign**
   - Packages not in any sync database
   - Usually AUR or manually built packages
   - Useful for tracking non-official packages

5. **All**
   - Every installed package
   - No filtering applied

### Switching Presets
- `Tab` - Cycle through presets in order
- `:preset <name>` - Jump directly to a preset

### Preset Indicator
- Status line shows current preset: `[Explicit] | 1,234 packages`

---

## Status Line

The status line appears at the bottom of the screen and provides contextual information.

### Normal Mode Display
```
[Explicit] | Showing 1-20 of 1,234 | Sort: Name ↑ | Filter: "firefox"
```

Components:
- Current preset/view name
- Visible row range and total count
- Active sort column and direction (if any)
- Active filter term (if any)

### Command Mode Display
```
:g 100_
```
Shows the command buffer with cursor position.

### Filter Mode Display
```
/firefox_
```
Shows the filter buffer with cursor position.

### Loading State
```
Loading packages...
```

### Error State
```
Error: Failed to load packages | Press ? for help
```

---

## Error Handling

### User Input Errors
- Invalid command: Show message "Unknown command: <input>"
- Invalid line number: "Line number out of range"
- Invalid sort column: "Unknown column: <input>"
- Invalid preset: "Unknown preset: <input>"

### System Errors
- Database error: Show error message, offer `:refresh` to retry
- Permission error: Explain issue and suggest solution
- Network error: Show error, indicate offline mode

### Error Display
- Errors appear in status line with red/error styling
- Press `Esc` or any key to dismiss error message

---

## Help Screen

📋 **Status**: Planned

Activated by pressing `?` in normal mode.

**Display**:
- Modal overlay with keybinding reference
- Organized by category (Navigation, Actions, Modes, etc.)
- Shows currently active keybindings from config

**Help Screen Controls**:
- `Esc` / `q` / `?` - Close help screen
- `↑` / `↓` - Scroll help content
- `/` - Search help text

---

## Package Management Actions

### Installation
📋 **Status**: Phase 4 - Planned

**Trigger**: `i` key or `:install` command

**Flow**:
1. Select package (from search results or info view)
2. Press `i` or `:install`
3. Confirmation prompt appears:
   ```
   Install: firefox 122.0-1
   Download size: 62.5 MB
   Installed size: 245 MB
   Dependencies: 5 packages

   Proceed? [Y/n]
   ```
4. On confirmation, show real-time progress
5. On completion, refresh package list and show message

### Removal
📋 **Status**: Phase 4 - Planned

**Trigger**: `r` key or `:remove` command

**Flow**:
1. Select installed package
2. Press `r` or `:remove`
3. Confirmation prompt appears:
   ```
   Remove: firefox 122.0-1
   Freed space: 245 MB
   Required by: 0 packages

   Proceed? [Y/n]
   ```
4. If reverse dependencies exist, show warning
5. On confirmation, show real-time progress
6. On completion, refresh package list and show message

### Safe Mode
- Prevent removal of critical packages (base, kernel, pacman)
- Show warning for packages with many reverse dependencies
- Offer cascade removal option (remove with dependencies)

---

## Search

### Local Filter (Filter Mode)
🚧 **Status**: Phase 1 - In Progress

- Instant, in-memory filtering
- Searches package names and descriptions
- Updates as you type
- Case-insensitive by default

### Remote Search
📋 **Status**: Phase 4 - Planned

**Command**: `:search <term>` or `:s <term>`

**Behavior**:
- Queries sync databases (official repositories)
- Shows results in main table view
- Indicates installation status per package
- Allows installation from search results (press `i`)
- Shows search result count in status line

**Search Result Columns**:
- Name, Version, Repository, Installed (✓/✗), Size, Description

### AUR Search
📋 **Status**: Future Enhancement

**Command**: `:aur <term>`

Similar to remote search but queries AUR database.

---

## Batch Operations

📋 **Status**: Future Enhancement

### Visual Mode (Multi-Select)
- `v` - Enter visual mode
- `↑` / `↓` - Extend selection
- `Space` - Toggle selection of current row
- `Esc` - Exit visual mode

### Batch Actions
- `:install` - Install all selected packages
- `:remove` - Remove all selected packages
- `:export` - Export selected package list to file

---

## Mouse Support

📋 **Status**: Future Enhancement (Optional)

When enabled in config:
- Click on row to select
- Click on column header to sort
- Scroll wheel to navigate
- Double-click to show info view

**Note**: Mouse support is optional and disabled by default to maintain keyboard-centric workflow.

---

## Configuration

All keybindings can be customized in `~/.config/pacviz/config.toml`:

```toml
[keybindings]
quit = ["q", "ctrl+c"]
filter = ["/"]
command = [":"]
info = ["enter", "i"]
help = ["?"]
sort_toggle = ["space"]
install = ["i"]
remove = ["r"]
```

---

## Implementation Status Summary

### ✅ Completed (Phase 0)
- Basic vertical navigation
- Page navigation (Ctrl+U/D)
- Jump navigation (g/G, home/end)
- Horizontal column navigation
- Quit functionality
- Window resizing
- Viewport scrolling and selection

### 🚧 In Progress (Phase 1)
- Command mode infrastructure
- Filter mode infrastructure
- Buffer management
- Status line mode display

### 📋 Planned

**Phase 2: Command System**
- Navigation commands (`:g`, `:t`, `:e`)
- Command parser enhancement
- Command execution

**Phase 3: Advanced Features**
- Sort toggle (Space)
- Preset/view switching (Tab)
- Info view (Enter/i)
- View reset after filter

**Phase 4: Polish**
- Additional commands (`:q`, `:sort`, `:preset`, `:filter`)
- Enhanced status line
- Command history
- Command autocomplete

**Future Enhancements**
- Package installation/removal
- Remote search
- AUR integration
- Batch operations
- Help screen
- Mouse support
- Plugin system
