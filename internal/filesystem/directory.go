package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type FileEntry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
	Mode    os.FileMode
}

type SortBy int

const (
	SortByName SortBy = iota
	SortBySize
	SortByModTime
)

func ListDirectory(path string) ([]FileEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	var files []FileEntry
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue // Skip entries we can't stat
		}

		files = append(files, FileEntry{
			Name:    e.Name(),
			Path:    filepath.Join(path, e.Name()),
			IsDir:   e.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Mode:    info.Mode(),
		})
	}

	// Sort: Directories first, then alphabetically
	SortEntries(files, SortByName)

	return files, nil
}

func SortEntries(entries []FileEntry, sortBy SortBy) {
	sort.Slice(entries, func(i, j int) bool {
		// Directories always come first
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}

		switch sortBy {
		case SortBySize:
			if entries[i].Size == entries[j].Size {
				return entries[i].Name < entries[j].Name
			}
			return entries[i].Size > entries[j].Size // Descending
		case SortByModTime:
			if entries[i].ModTime == entries[j].ModTime {
				return entries[i].Name < entries[j].Name
			}
			return entries[i].ModTime.After(entries[j].ModTime) // Newest first
		default: // SortByName
			return entries[i].Name < entries[j].Name
		}
	})
}
