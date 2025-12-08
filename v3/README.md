# PacViz v3

A terminal user interface (TUI) for browsing and managing Arch Linux packages.

## Phase 1 - Complete ✓

Phase 1 implements a minimal viable product that displays all installed packages in a table format.

### Features Implemented

- **Package Loading**: Loads all installed packages from the local pacman database
- **Table Display**: Shows packages in a formatted table with columns:
  - Name
  - Version
  - Size (human-readable format)
  - Install Date
  - Description
- **Header Row**: Column names with sort indicators
- **Styled Rendering**:
  - Alternating row colors for readability
  - Selected row highlighting
  - Custom color scheme
  - Horizontal cell padding (1 character on each side)
  - Aligned headers and cells
- **Responsive Layout**: Column widths calculated based on terminal size
  - **Name**: 40 characters (fixed, 2x other columns)
  - **Version**: 20 characters (fixed, fits "20250814.1-1")
  - **Size**: 10 characters (fixed, fits "999.9 MB")
  - **Install Date**: 12 characters (fixed, fits "2025-11-30")
  - **Description**: Auto-sized (grows to fill remaining terminal width, min 20 chars)
  - Table fills full terminal width (total fixed: 82 chars + description)

### Architecture

```
v3/
├── cmd/pacviz/          # Application entry point
├── internal/
│   ├── app/             # Bubble Tea MVU (Model-View-Update)
│   ├── domain/          # Business logic & data structures
│   ├── repository/      # Package data access (ALPM)
│   ├── ui/
│   │   ├── column/      # Column definitions & width calculation
│   │   ├── renderer/    # UI rendering components
│   │   ├── styles/      # Colors & themes
│   │   └── viewport/    # Viewport state & logic
│   ├── command/         # Command system (stubbed)
│   └── config/          # Configuration (stubbed)
└── pkg/version/         # Semantic version comparison (stubbed)
```

### Building

```bash
go build ./cmd/pacviz
```

### Running

```bash
./pacviz
```

**Note**: Requires running on Arch Linux or compatible system with pacman installed.

### Testing Data Loading

The package loading was verified successfully:
- Loaded 912 packages from test system
- Successfully converted to table rows
- All package metadata correctly parsed

### Current Status

**Working**:
- ✅ ALPM repository integration
- ✅ Package loading (all installed packages)
- ✅ Package to row conversion
- ✅ Column width calculation algorithm
- ✅ Header rendering with sort indicators
- ✅ Table rendering with alternating colors
- ✅ Viewport state management
- ✅ Full UI composition

**Not Yet Implemented** (planned for later phases):
- Navigation (arrow keys, page up/down)
- Sorting (by column)
- Filtering (text search)
- Command mode (`:` commands)
- Package installation/removal
- Info view (detailed package information)
- Preset views (explicit, orphans, foreign)
- Scrolling for large package lists

## Next Steps

Phase 2 will add:
1. Keyboard navigation (up/down arrows)
2. Scrolling for viewing all packages
3. Column selection (left/right arrows)
4. Sorting by column (spacebar to toggle)

See `DESIGN.md` for the complete feature roadmap.

## Project Statistics

- **Lines of Code**: 36 Go files
- **Binary Size**: 4.6 MB (with debug symbols)
- **Dependencies**:
  - Bubble Tea (TUI framework)
  - Lipgloss (styling)
  - go-alpm (pacman interface)
  - go-pacmanconf (pacman config parser)
