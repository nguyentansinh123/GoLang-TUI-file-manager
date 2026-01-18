package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tansinhnguyen123/my-tui-explorer/internal/filesystem"
)

type Explorer struct {
	app           *tview.Application
	list          *tview.List
	parentList    *tview.List
	previewPane   *tview.TextView
	statusBar     *tview.TextView
	pathBar       *tview.TextView
	currentPath   string
	entries       []filesystem.FileEntry
	parentEntries []filesystem.FileEntry
	showHidden    bool
	searchMode    bool
	searchQuery   string
	selectedIndex int
	previewPath   string
}

func NewExplorer() *Explorer {
	return &Explorer{
		app:        tview.NewApplication(),
		showHidden: false,
	}
}

func (e *Explorer) Run() error {
	e.parentList = tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDefault)
	e.parentList.SetBorder(true).SetTitle("┌─ Parent ─").SetBorderPadding(0, 0, 1, 0)

	e.list = tview.NewList().
		ShowSecondaryText(true).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkGreen)
	e.list.SetBorder(true).SetTitle("┌─ Current ─").SetBorderPadding(0, 0, 1, 0)

	e.previewPane = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignLeft)
	e.previewPane.SetBorder(true).SetTitle("┌─ Preview ─").SetBorderPadding(0, 0, 1, 1)
	e.previewPane.SetText(" [dim]Select a file to preview[white]")

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	header.SetText("[cyan]╭───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮\n│[yellow]  TUI File Explorer[cyan] - [white]Modern File Browser[cyan]                                                                                 │\n╰───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯[white]")

	e.pathBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	e.pathBar.SetText("")

	e.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	e.updateStatusBar()

	e.setupKeybindings()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/"
	}
	e.navigate(homeDir)

	return e.app.SetRoot(e.createLayout(), true).EnableMouse(true).Run()
}

func (e *Explorer) createLayout() *tview.Flex {
	// 3-pane layout: parent | current | preview
	mainFlex := tview.NewFlex().
		AddItem(e.parentList, 0, 1, false).
		AddItem(e.list, 0, 2, true).
		AddItem(e.previewPane, 0, 1, false)

	// Top section with header
	headerFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText(" "), 0, 0, false)

	// Full layout
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 3, 0, false).
		AddItem(e.pathBar, 2, 0, false).
		AddItem(mainFlex, 0, 1, true).
		AddItem(e.statusBar, 2, 0, false)
}

func (e *Explorer) updateStatusBar() {
	if e.searchMode {
		e.statusBar.SetText(fmt.Sprintf(" [yellow]🔍 Search: [white]%-30s [yellow]│[white] Press [yellow]Enter[white] to search, [yellow]Esc[white] to cancel", e.searchQuery))
	} else {
		hiddenStr := ""
		if e.showHidden {
			hiddenStr = " [green]·[white] Hidden files [green]ON"
		}

		e.statusBar.SetText(fmt.Sprintf(" [cyan]↑[white]/[cyan]↓[white]: navigate  │  [cyan]Enter[white]/[cyan]l[white]: open  │  [cyan]h[white]: back  │  [cyan]~[white]: home  │  [cyan].[white]: toggle hidden  │  [cyan]?[white]: help  │  [cyan]q[white]: quit%s", hiddenStr))
	}
}

func (e *Explorer) setupKeybindings() {
	e.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q': // Quit
				e.app.Stop()
				return nil
			case 'j': // Move down
				current := e.list.GetCurrentItem()
				if current < e.list.GetItemCount()-1 {
					e.list.SetCurrentItem(current + 1)
					e.updateListPreview()
				}
				return nil
			case 'k': // Move up
				current := e.list.GetCurrentItem()
				if current > 0 {
					e.list.SetCurrentItem(current - 1)
					e.updateListPreview()
				}
				return nil
			case 'l': // Enter directory
				e.openSelected()
				return nil
			case 'h': // Go to parent directory
				e.goBack()
				return nil
			case 'g': // Go to top
				e.list.SetCurrentItem(0)
				e.updateListPreview()
				return nil
			case 'G': // Go to bottom
				e.list.SetCurrentItem(e.list.GetItemCount() - 1)
				e.updateListPreview()
				return nil
			case '.': // Toggle hidden files
				e.showHidden = !e.showHidden
				e.refresh()
				return nil
			case '~': // Go to home
				homeDir, _ := os.UserHomeDir()
				e.navigate(homeDir)
				return nil
			case '/': // Go to root
				e.navigate("/")
				return nil
			case '?': // Show help
				e.showHelp()
				return nil
			case 's': // Start search
				e.startSearch()
				return nil
			}
		case tcell.KeyEnter: // Enter directory or open file
			e.openSelected()
			return nil
		case tcell.KeyEscape: // Go back
			e.goBack()
			return nil
		}
		return event
	})
}

func (e *Explorer) navigate(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	entries, err := filesystem.ListDirectory(absPath)
	if err != nil {
		e.statusBar.SetText(fmt.Sprintf(" [red]✗ Error: %v", err))
		return
	}

	e.currentPath = absPath

	// Load parent directory
	parentPath := filepath.Dir(absPath)
	if parentPath != absPath {
		parentEntries, _ := filesystem.ListDirectory(parentPath)
		e.parentEntries = parentEntries
	}

	e.entries = entries
	e.selectedIndex = 0
	e.refresh()
}

func (e *Explorer) refresh() {
	// Update path bar with better styling
	homeDir, _ := os.UserHomeDir()
	displayPath := e.currentPath
	if strings.HasPrefix(displayPath, homeDir) {
		displayPath = "~" + strings.TrimPrefix(displayPath, homeDir)
	}
	e.pathBar.SetText(fmt.Sprintf(" [yellow]📂  [white]%s", displayPath))

	// Refresh parent list
	e.parentList.Clear()
	for _, entry := range e.parentEntries {
		if !e.showHidden && len(entry.Name) > 0 && entry.Name[0] == '.' {
			continue
		}
		var icon string
		if entry.IsDir {
			icon = "📁"
		} else {
			icon = "📄"
		}
		text := fmt.Sprintf("%s  %s", icon, entry.Name)
		e.parentList.AddItem(text, "", 0, nil)
	}

	// Refresh main list
	e.list.Clear()
	for _, entry := range e.entries {
		// Skip hidden files if not showing
		if !e.showHidden && len(entry.Name) > 0 && entry.Name[0] == '.' {
			continue
		}

		var icon, color, nameText string
		var sizeStr string

		if entry.IsDir {
			icon = "📁"
			color = "[cyan]"
			nameText = entry.Name
			sizeStr = fmt.Sprintf("[dim]%d items", countDirItems(entry.Path))
		} else {
			icon = "📄"
			color = "[white]"
			nameText = entry.Name
			sizeStr = fmt.Sprintf("[dim]%s", formatSize(entry.Size))
		}

		mainText := fmt.Sprintf("%s%s  %s[white]", color, icon, nameText)
		e.list.AddItem(mainText, sizeStr, 0, nil)
	}

	// Update status bar
	e.updateStatusBar()

	// Update preview for selected item
	if e.selectedIndex >= 0 && e.selectedIndex < len(e.entries) {
		e.updatePreview(e.entries[e.selectedIndex])
	}
}

func (e *Explorer) openSelected() {
	idx := e.list.GetCurrentItem()

	// Filter entries same way as refresh
	var visibleEntries []filesystem.FileEntry
	for _, entry := range e.entries {
		if !e.showHidden && len(entry.Name) > 0 && entry.Name[0] == '.' {
			continue
		}
		visibleEntries = append(visibleEntries, entry)
	}

	if idx < 0 || idx >= len(visibleEntries) {
		return
	}

	entry := visibleEntries[idx]
	if entry.IsDir {
		e.navigate(entry.Path)
	} else {
		// Open file with default editor
		e.openFile(entry.Path)
	}
}

func (e *Explorer) openFile(path string) {
	// Stop the app, open file, then resume
	e.app.Suspend(func() {
		// Try different common editors
		editors := []string{"nano", "vim", "vi", "emacs", "gedit", "code"}

		for _, editor := range editors {
			cmd := exec.Command(editor, path)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err == nil {
				return
			}
		}

		// If no editor found, try xdg-open (Linux) or open (macOS)
		var cmd *exec.Cmd
		if os.Getenv("OSTYPE") == "darwin" {
			cmd = exec.Command("open", path)
		} else {
			cmd = exec.Command("xdg-open", path)
		}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	})
}

func (e *Explorer) updatePreview(entry filesystem.FileEntry) {
	e.previewPath = entry.Path

	if entry.IsDir {
		// Show directory info
		fileCount := countDirItems(entry.Path)
		info := fmt.Sprintf("[cyan]📁 Directory[white]\n\n[dim]Name:[white] %s\n[dim]Path:[white] %s\n[dim]Items:[white] %d\n\n[yellow]Press Enter to open", entry.Name, entry.Path, fileCount)
		e.previewPane.SetText(info)
	} else {
		// Show file info with preview
		info := fmt.Sprintf("[cyan]📄 File[white]\n\n[dim]Name:[white] %s\n[dim]Size:[white] %s\n[dim]Path:[white] %s\n\n", entry.Name, formatSize(entry.Size), entry.Path)

		// Try to preview file content if text file
		if isTextFile(entry.Name) {
			content := previewFile(entry.Path)
			info += "[cyan]━━━ Preview ━━━[white]\n" + content
		} else {
			info += "[dim](Binary file - no preview)[white]"
		}

		e.previewPane.SetText(info)
	}
}

func (e *Explorer) updateListPreview() {
	idx := e.list.GetCurrentItem()

	// Filter entries same way as refresh
	var visibleEntries []filesystem.FileEntry
	for _, entry := range e.entries {
		if !e.showHidden && len(entry.Name) > 0 && entry.Name[0] == '.' {
			continue
		}
		visibleEntries = append(visibleEntries, entry)
	}

	if idx >= 0 && idx < len(visibleEntries) {
		e.updatePreview(visibleEntries[idx])
	}
}

func countDirItems(path string) int {
	entries, err := filesystem.ListDirectory(path)
	if err != nil {
		return 0
	}
	return len(entries)
}

func isTextFile(filename string) bool {
	textExtensions := map[string]bool{
		".txt": true, ".md": true, ".json": true, ".yaml": true, ".yml": true,
		".go": true, ".py": true, ".js": true, ".ts": true, ".css": true,
		".html": true, ".xml": true, ".csv": true, ".sh": true, ".bash": true,
		".conf": true, ".config": true, ".toml": true, ".ini": true,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return textExtensions[ext]
}

func previewFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return "[red]Error reading file[white]"
	}

	text := string(content)
	if len(text) > 500 {
		text = text[:500] + "[dim]\n... (truncated)[white]"
	}

	return "[white]" + text
}

func (e *Explorer) goBack() {
	parent := filepath.Dir(e.currentPath)
	if parent != e.currentPath {
		e.navigate(parent)
	}
}

func (e *Explorer) showHelp() {
	help := `[yellow]━━━ TUI Explorer Help ━━━[white]

[cyan]Navigation:[white]
  [yellow]j/k[white] or [yellow]↑/↓[white]  Navigate up/down
  [yellow]h[white] or [yellow]Esc[white]    Go to parent directory
  [yellow]Enter/l[white]         Open directory/file
  [yellow]~[white]                Go to home directory
  [yellow]/[white]                Go to root directory

[cyan]View:[white]
  [yellow]g[white]                Go to top
  [yellow]G[white]                Go to bottom
  [yellow].[white]                Toggle hidden files
  [yellow]s[white]                Search files
  [yellow]?[white]                Show this help

[cyan]Other:[white]
  [yellow]q[white]                Quit application

[yellow]Press any key to close...[white]`

	modal := tview.NewModal().
		SetText(help).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			e.app.SetRoot(e.createLayout(), true)
		})

	e.app.SetRoot(modal, true)
}

func (e *Explorer) startSearch() {
	inputField := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(30)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			query := inputField.GetText()
			if query != "" {
				e.filterSearch(query)
			}
			e.app.SetRoot(e.createLayout(), true)
		} else if key == tcell.KeyEscape {
			e.app.SetRoot(e.createLayout(), true)
		}
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(inputField, 3, 0, true).
		AddItem(tview.NewTextView().SetText(" [yellow]Enter[white]: search | [yellow]Esc[white]: cancel"), 1, 0, false).
		AddItem(nil, 0, 1, false)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(flex, 40, 0, true).
		AddItem(nil, 0, 1, false)

	e.app.SetRoot(modal, true).SetFocus(inputField)
}

func (e *Explorer) filterSearch(query string) {
	query = strings.ToLower(query)
	e.list.Clear()
	count := 0

	for _, entry := range e.entries {
		if !e.showHidden && len(entry.Name) > 0 && entry.Name[0] == '.' {
			continue
		}

		if strings.Contains(strings.ToLower(entry.Name), query) {
			var icon, color string
			var sizeStr string

			if entry.IsDir {
				icon = "📁"
				color = "blue"
				sizeStr = "<DIR>"
			} else {
				icon = "📄"
				color = "white"
				sizeStr = formatSize(entry.Size)
			}

			mainText := fmt.Sprintf("[%s]%s %s[white]", color, icon, entry.Name)
			secondaryText := fmt.Sprintf("    %s", sizeStr)
			e.list.AddItem(mainText, secondaryText, 0, nil)
			count++
		}
	}

	if count == 0 {
		e.statusBar.SetText(" [red]No matches found")
	} else {
		e.statusBar.SetText(fmt.Sprintf(" [green]Found %d matches", count))
	}
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
