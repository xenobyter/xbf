# xbf

Keyboard-driven dual-pane file manager for the terminal (Go + tview/tcell).

`xbf` provides two independent directory views (left/right) and supports common file operations directly from the keyboard.

## Features

- Dual-pane navigation
- Per-directory sorting:
  - directories first
  - then alphabetical (case-insensitive, original casing as tie-breaker)
- File info in the footer (name + formatted size)
- Fullscreen file preview page
  - readable text files as text
  - binary files as hex dump
  - syntax highlighting for text previews (lexer-based)
- Editor-ready document model for text files
  - open current file in an edit session
  - save with `Ctrl+S`
  - close with `Esc`
- Simple Hex editor for binary files
  - open current file in a hex edit session
  - save with `Ctrl+S`
  - close with `Esc`
- Multi-selection per pane
- File operations:
  - copy
  - move
  - rename
  - delete (with confirmation dialog)
- Recursive directory handling
- Symlink support for copy/move/rename/delete
- i18n (German/English), automatic language detection via `LC_ALL`/`LANG` with English fallback
- Comprehensive unit tests for filesystem operations, input mapping, selection logic, and i18n

## Requirements

- Go (version defined in `go.mod`)

## Build & Start

```bash
go build -o xbf .
./xbf
```

Or run directly:

```bash
go run .
```

## Keyboard Shortcuts

| Key | Action |
| --- | --- |
| `Tab` | Switch active pane |
| `→` | Enter selected directory |
| `←` | Go to parent directory |
| `s` / `Space` | Select/deselect item |
| `c` | Copy selection (or current item) to inactive pane |
| `m` | Move selection (or current item) to inactive pane |
| `r` | Rename current item (input dialog) |
| `d` / `Del` | Delete selection (or current item), with confirmation |
| `n` | Create a new empty file |
| `N` | Create a new directory |
| `p` | Open fullscreen preview for current file |
| `e` | Open current file in editor session |
| `q` / `Esc` | Quit application |

Preview page:

- `Esc` = close preview
- `↓` = scroll down
- `↑` = scroll up
- `PgDn` / `PgUp` = page scroll
- `Home` / `End` = start/end

Delete confirmation dialog:

- `y` / `j` = confirm
- `n` / `Esc` = cancel

## UI Layout

- **Header**: current working directory of the active pane
- **Left pane**: left directory contents
- **Right pane**: right directory contents
- **Footer**: file info, success messages, error messages

Icons:

- 📁 Directory
- 📄 File

## Internationalization

Supported languages:

- German (`de`)
- English (`en`)

Examples:

```bash
LANG=de_DE.UTF-8 ./xbf
LANG=en_US.UTF-8 ./xbf
```

## Tests

```bash
go test ./...
```

Optional verbose output:

```bash
go test -v ./...
```

## License

This project is released under the Unlicense. See `LICENSE` for details.
