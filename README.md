# pacviz

A terminal user interface for managing Arch Linux packages via pacman.

![screenshot](screenshot1.png)

## Features

Pacviz displays your installed packages in a table, saving you from running `-Qe` all the time. Switch between various presets such as Explicit or Orphans. Install and remove packages from the official Arch repository and the AUR (using either `paru` or `yay`).

## Usage

```bash
pacviz
pacviz -c config.toml
```

### Keybinds

- `j/k` or `↑/↓`: Navigate
- `/`: Filter packages
- `tab`: Cycle presets
- `enter`: Toggle detail panel
- `:`: Command mode
- `q`: Quit

### Commands

- `:install` or `:i`: Install selected package
- `:remove` or `:r`: Remove selected package
- `:sort <col> <asc|desc>`: Sort by column
- `:filter <term>`: Apply filter
- `:preset <name>`: Switch preset
- `:help` or `:?`: Show help
