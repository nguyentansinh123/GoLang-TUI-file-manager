# GoLang TUI File Manager

A modern, feature-rich **Terminal User Interface (TUI)** file explorer built with Go.  Navigate your filesystem with ease using an intuitive and responsive interface powered by `tview` and `tcell`.

![Language](https://img.shields.io/badge/language-Go-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

### Core Functionality
- 📁 **Dual-pane Navigation** - View parent and current directory simultaneously
- 👁️ **Live File Preview** - Preview text files and directory contents in real-time
- 🔍 **Smart Search** - Find files and directories quickly with fuzzy search
- 👻 **Hidden Files Toggle** - Show/hide dotfiles with a single keystroke
- 🎨 **Syntax Highlighting** - Color-coded file types with emoji icons
- 📊 **File Information** - Display file size, modification time, and permissions
- ⌨️ **Vim-like Keybindings** - Navigate efficiently with familiar shortcuts

### Advanced Features
- Open files with your preferred editor (nano, vim, vi, emacs, code, etc.)
- Automatic directory sorting (directories first, then alphabetical)
- Cross-platform support (Linux, macOS, Windows)
- Beautiful border styling and status bar
- Mouse support for point-and-click navigation

## 🚀 Quick Start

### Prerequisites
- Go 1.25.5 or higher
- A terminal with color support

### Installation

1. **Clone the repository:**
```bash
git clone https://github.com/nguyentansinh123/GoLang-TUI-file-manager. git
cd GoLang-TUI-file-manager
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Build the application:**
```bash
go build -o explorer cmd/explorer/main.go
```

4. **Run the explorer:**
```bash
./explorer
```

## ⌨️ Keybindings

### Navigation
| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `h` / `Esc` | Go to parent directory |
| `l` / `Enter` | Open directory or file |
| `g` | Go to top of list |
| `G` | Go to bottom of list |
| `~` | Go to home directory |
| `/` | Go to root directory |

### View & Search
| Key | Action |
|-----|--------|
| `.` | Toggle hidden files |
| `s` | Search files |
| `?` | Show help menu |

### Other
| Key | Action |
|-----|--------|
| `q` | Quit application |

## 🎯 Usage Examples

### Basic Navigation
1. Launch the explorer with `./explorer`
2. Use `j`/`k` or arrow keys to navigate
3. Press `Enter` or `l` to open a directory
4. Press `h` or `Esc` to go back

### Search for Files
1. Press `s` to activate search mode
2. Type your search query
3. Press `Enter` to filter results
4. Press `Esc` to cancel

### Toggle Hidden Files
Press `.` to show or hide dotfiles (files starting with `.`)

## 📁 Project Structure

```
GoLang-TUI-file-manager/
├── cmd/
│   └── explorer/
│       └── main.go           # Application entry point
├── internal/
│   ├── filesystem/
│   │   └── directory. go      # File system operations & sorting
│   └── ui/
│       └── ui.go             # TUI interface & keybindings
├── go.mod                    # Go module dependencies
├── go.sum                    # Dependency checksums
└── README.md
```

## 🛠️ Dependencies

This project uses the following Go libraries: 

- **[tview](https://github.com/rivo/tview)** v0.42.0 - Terminal UI framework
- **[tcell](https://github.com/gdamore/tcell)** v2.13.7 - Terminal cell-based display
- **[gopsutil](https://github.com/shirou/gopsutil)** v3.24.5 - System utilities

## 🎨 UI Components

### Three-Pane Layout
1. **Parent Directory** (Left) - Shows contents of the parent folder
2. **Current Directory** (Center) - Main navigation pane with file list
3. **Preview Pane** (Right) - File preview and directory information

### Status Bar
- Displays current path with home directory aliasing
- Shows file count and navigation information

## 🔧 Configuration

The explorer supports various text editors for opening files:
- `nano` (default fallback)
- `vim`
- `vi`
- `emacs`
- `gedit`
- `code` (VS Code)
- System default (`xdg-open` on Linux, `open` on macOS)

## 🐛 Troubleshooting

### Editor Not Opening Files
The application tries common editors in sequence.  If none work:
1. Install a supported editor (e.g., `sudo apt install nano`)
2. Or ensure `xdg-open` (Linux) / `open` (macOS) is configured

### Display Issues
- Ensure your terminal supports 256 colors
- Try resizing the terminal window
- Check that your terminal emulator is up to date

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Future Enhancements

- [ ] File operations (copy, move, delete, rename)
- [ ] Multiple sorting options (size, date, type)
- [ ] Custom color themes
- [ ] File permissions editing
- [ ] Archive preview (zip, tar, etc.)
- [ ] Bookmarks for quick directory access
- [ ] Split screen mode

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👤 Author

**Nguyen Tan Sinh**
- GitHub: [@nguyentansinh123](https://github.com/nguyentansinh123)

## 🌟 Acknowledgments

- Built with [tview](https://github.com/rivo/tview) by rivo
- Terminal handling by [tcell](https://github.com/gdamore/tcell)
- Inspired by modern terminal file managers like `ranger` and `lf`

---

⭐ If you find this project useful, please consider giving it a star! 
