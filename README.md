# xbf - Dual Pane Terminal File Manager

A lightweight dual-pane file manager for the terminal, written in Go and powered by the **tview** and **tcell** libraries.

The application provides two independent directory views, allowing users to navigate the filesystem efficiently using only the keyboard. It also includes basic internationalization (English and German) and a clean terminal-based user interface.

## Features

- Dual-pane file browser
- Independent navigation in each pane
- Keyboard-driven interface
- File and directory icons
- Automatic file sorting:
  - Directories first
  - Alphabetical order (case-insensitive)
- Automatic language detection via environment variables
- English and German translations
- Unit tests for:
  - File sorting
  - Internationalization

## Dependencies

- [tview](https://github.com/rivo/tview) – Terminal UI framework
- [tcell](https://github.com/gdamore/tcell) – Terminal event handling

## Keyboard Shortcuts

| Key | Action |
| ------- | -------- |
| `Tab` | Switch between left and right pane |
| `→` | Enter selected directory |
| `←` | Navigate to parent directory |
| `s`/`Space` | Select/deselect item |
| `q`/`Esc` | Quit application |

## User Interface

The interface consists of four sections:

### Header

Displays the current working directory of the active pane.

### Left Pane

Shows the contents of the left working directory.

### Right Pane

Shows the contents of the right working directory.

### Footer

Reserved for status messages and future enhancements.

## File and Directory Icons

The application uses Unicode icons to distinguish between files and directories.

| Icon | Description |
| ------- | ------------- |
| 📁 | Directory |
| 📄 | File |

## Sorting Behavior

Entries are sorted according to the following rules:

1. Directories before files
2. Alphabetical order ignoring case
3. Original letter casing used as a tie-breaker

Example:

```text
📁 alpha_dir
📁 Beta_dir
📄 a.txt
📄 apple.txt
📄 B.txt
📄 zebra.txt
```

## Internationalization

The application automatically detects the system language using the `LC_ALL` and `LANG` environment variables.

Currently supported languages:

- English (`en`)
- German (`de`)

Examples:

```bash
LANG=en_US.UTF-8 ./filemanager
```

```bash
LANG=de_DE.UTF-8 ./filemanager
```

If no supported language is detected, the application falls back to English.

## Project Structure

```text
.
├── main.go
├── app.go            # Application initialization and layout
├── ui.go             # Directory navigation and list handling
├── input.go          # Keyboard input processing
├── file.go           # Filesystem utilities
├── i18n.go           # Internationalization support
├── file_test.go      # Sorting tests
├── i18n_test.go      # Translation tests
└── go.mod
```

## Architecture

### App

The `App` structure acts as the central application controller and manages:

- tview application instance
- Header and footer
- Left and right directory panes
- Current working directories
- Directory contents
- Language management

### FileInfo

Represents a file or directory entry:

```go
type FileInfo struct {
    Name  string
    Path  string
    IsDir bool
    Size  int64
}
```

### I18n

Handles language detection and translation lookup.

### Event Handling

Keyboard events are normalized into symbolic actions and dispatched through a central action map.

```text
q / Esc  → quit
Tab      → switch pane
→        → enter directory
←        → parent directory
```

## Running Tests

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

## Error Handling

The application gracefully handles common filesystem errors such as:

- Unable to determine the current working directory
- Unable to read directory contents
- Invalid navigation attempts

Error messages are localized according to the selected language.

## Future Enhancements

Potential future features include:

- Copy files and directories
- Move/rename files
- Delete operations
- Create directories
- File size display
- Status bar information
- Configurable themes
- File preview panel

## License

This project is public domain under the The Unlicense. See the `LICENSE` file for details.
